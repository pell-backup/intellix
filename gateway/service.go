package taskgateway

import (
	"context"
	"cosmossdk.io/math"
	"encoding/json"
	"errors"
	"fmt"
	dvslog "github.com/0xPellNetwork/pelldvs/libs/log"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"net"
	"net/rpc"
	"os"
	"time"

	dataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"sync"
)

type TaskGateway struct {
	service.BaseService

	server     *rpc.Server
	serverAddr string
	listener   net.Listener

	cfg *TaskGatewayCfg
	ctx context.Context

	logger     dvslog.Logger
	ethClient  *ethclient.Client
	privateKey *keystore.Key

	contractDataOracle *dataOracle.ContractDataOracle
	taskMap            sync.Map
	nonceMap           sync.Map
}

// TODO: put it in a common location
func packUint256(value *big.Int) ([]byte, error) {
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}

	arguments := abi.Arguments{{Type: uint256Type}}

	return arguments.Pack(value)
}

func NewTaskGateway(logger dvslog.Logger, ctx context.Context, cfg *TaskGatewayCfg) (*TaskGateway, error) {
	ethClient, err := ethclient.Dial(cfg.EthEndpoint)
	if err != nil {
		return nil, err
	}

	contract, err := dataOracle.NewContractDataOracle(common.HexToAddress(cfg.ContractAddress), ethClient)
	if err != nil {
		return nil, err
	}

	// Read private key
	keyJSON, err := os.ReadFile(cfg.PrivateKeyStorePath)
	if err != nil {
		logger.Error("Failed to read private key file", "error", err)
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	key, err := keystore.DecryptKey(keyJSON, "")
	if err != nil {
		logger.Error("Failed to decrypt private key", "error", err)
		return nil, fmt.Errorf("failed to decrypt private key: %v", err)
	}

	server := rpc.NewServer()
	tg := &TaskGateway{
		server:             server,
		cfg:                cfg,
		ctx:                ctx,
		logger:             logger,
		ethClient:          ethClient,
		privateKey:         key,
		serverAddr:         cfg.ServerAddr,
		contractDataOracle: contract,
		taskMap:            sync.Map{},
		nonceMap:           sync.Map{},
	}
	if err := server.Register(tg); err != nil {
		logger.Error("Failed to register RPC server", "error", err)
		return nil, fmt.Errorf("failed to register RPC server: %v", err)
	}

	tg.BaseService = *service.NewBaseService(nil, "TaskGateway", tg)
	return tg, nil
}

func (tg *TaskGateway) OnStart() error {
	var err error
	tg.listener, err = net.Listen("tcp", tg.serverAddr)
	if err != nil {
		tg.logger.Error("Failed to start listener", "address", tg.serverAddr, "error", err)
		panic(err)
	}

	go tg.server.Accept(tg.listener)
	return nil
}

func (tg *TaskGateway) OnStop() {
	tg.ethClient.Close()
	if tg.listener != nil {
		tg.listener.Close()
	}
}

func (tg *TaskGateway) RespondToTask(req *RPCVoteFinalizedRequestIn, resp *RespondToTaskResponse) error {
	err := tg.handleResponse(context.Background(), req)
	if err != nil {
		resp.Error = err.Error()
		return err
	}
	resp.Error = ""
	return nil
}

func (tg *TaskGateway) getAuthOpts(chainId int64) (*bind.TransactOpts, error) {
	// Create transaction authenticator
	auth, err := bind.NewKeyedTransactorWithChainID(tg.privateKey.PrivateKey, big.NewInt(chainId))
	if err != nil {
		tg.logger.Error("Failed to create transaction authenticator", "error", err)
		return nil, fmt.Errorf("failed to create transaction authenticator: %v", err)
	}

	return auth, nil
}

func (tg *TaskGateway) handleResponse(ctx context.Context, response *RPCVoteFinalizedRequestIn) error {
	value, loaded := tg.taskMap.LoadOrStore(response.TaskRaw.TaskIndex, response)
	if loaded {
		existingResponse := value.(*RPCVoteFinalizedRequestIn)
		if tg.shouldReplaceResponse(ctx, existingResponse, response) {
			tg.taskMap.Store(response.TaskRaw.RequestID, response)
			return tg.wrapSubmitToChain(ctx, response)
		}
	} else {
		return tg.wrapSubmitToChain(ctx, response)
	}
	return nil
}

func (tg *TaskGateway) shouldReplaceResponse(ctx context.Context, existing, new *RPCVoteFinalizedRequestIn) bool {
	// TODO: add security threshold comparison
	return false
}

func (tg *TaskGateway) wrapSubmitToChain(ctx context.Context, request *RPCVoteFinalizedRequestIn) error {
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	var err error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				tg.logger.Error("Failed to submit vote finalized request", "error", r)
				err = fmt.Errorf("%v", r)
			}
			wg.Done()
		}()
		err = tg.submitToChain(ctx, request)
	}()
	wg.Wait()
	return err
}

func (tg *TaskGateway) submitToChain(ctx context.Context, response *RPCVoteFinalizedRequestIn) error {
	jsData, _ := json.Marshal(response)
	tg.logger.Info("TaskGateway.submitToChain request", "data", string(jsData))

	// Validate BLS signature components
	if err := validateBLSComponents(response.ValidatedData); err != nil {
		tg.logger.Error("Invalid BLS signature components", "error", err)
		return fmt.Errorf("invalid BLS components: %v", err)
	}

	feeTokenAddr, err := convertAddressToString(response.TaskRaw.FeeToken)
	if err != nil {
		tg.logger.Error("Error converting fee token address", "err", err)
		return err
	}
	cbAddr, err := convertAddressToString(response.TaskRaw.CallbackAddress)
	if err != nil {
		tg.logger.Error("Error converting callback address", "err", err)
		return err
	}

	paymentInt, ok := math.NewIntFromString(response.TaskRaw.Payment)
	if !ok {
		tg.logger.Error("Error converting taskRaw payment", "payment", response.TaskRaw.Payment)
		return fmt.Errorf("error converting taskRaw payment")
	}
	task := dataOracle.IDataOracleTask{
		TaskType:                 math.NewInt(response.TaskRaw.TaskType).BigInt(),
		RequestId:                [32]byte(response.TaskRaw.RequestID),
		FeeToken:                 *feeTokenAddr,
		AdvanceDecode:            response.TaskRaw.AdvanceDecode,
		Payment:                  paymentInt.BigInt(),
		RequestData:              response.TaskRaw.RequestData,
		CallbackAddress:          *cbAddr,
		CallbackFunctionId:       [4]byte(response.TaskRaw.CallbackFunctionID),
		TaskCreatedBlock:         response.TaskRaw.TaskCreatedBlock,
		GroupNumbers:             response.TaskRaw.QuorumNumbers,
		GroupThresholdPercentage: response.TaskRaw.QuorumThresholdPercentage,
	}

	sign := dataOracle.IBLSSignatureVerifierNonSignerStakesAndSignature{
		NonSignerGroupBitmapIndices: response.ValidatedData.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:            convertNonSignersPubkeysG1(response.ValidatedData.NonSignersPubkeysG1),
		GroupApks:                   convertQuorumApks(response.ValidatedData.QuorumApksG1),
		ApkG2:                       convertApkG2(response.ValidatedData.SignersApkG2),
		Sigma:                       convertSigma(response.ValidatedData.SignersAggSigG1),
		GroupApkIndices:             response.ValidatedData.QuorumApkIndices,
		TotalStakeIndices:           response.ValidatedData.TotalStakeIndices,
		NonSignerStakeIndices:       response.ValidatedData.NonSignerStakeIndices,
	}

	taskResp := dataOracle.IDataOracleTaskResponse{
		ReferenceTaskIndex: response.TaskRaw.TaskIndex,
		Data:               response.RespToTaskData,
	}

	authOpts, err := tg.getAuthOpts(response.ChainID)
	if err != nil {
		return err
	}

	// For debug purpose,
	// set gas limit to 1000000 to bypass transaction pre-execution and force broadcast
	// authOpts.GasLimit = 1000000

	transaction, err := tg.contractDataOracle.ResponseToTask(authOpts, task, taskResp, sign)
	if err != nil {
		// Try to get the failed transaction receipt
		tg.logger.Error("Error assembling RequestPrice tx",
			"err", err,
			"task", fmt.Sprintf("%+v", task),
			"taskResp", fmt.Sprintf("%+v", taskResp),
			"sign", fmt.Sprintf("%+v", sign))
		return err
	}

	if transaction != nil {
		_ = tg.queryTransaction(ctx, transaction)
	}

	return nil
}

func (tg *TaskGateway) queryTransaction(ctx context.Context, tx *types.Transaction) error {
	tg.logger.Info("Transaction submitted", "txHash", tx.Hash().Hex())

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	receipt, err := bind.WaitMined(timeoutCtx, tg.ethClient, tx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			tg.logger.Error("Transaction confirmation timeout",
				"txHash", tx.Hash().Hex())
			return fmt.Errorf("transaction confirmation timeout: %s", tx.Hash().Hex())
		}
		tg.logger.Error("Error waiting for transaction to be mined",
			"txHash", tx.Hash().Hex(),
			"error", err)
		return err
	}

	tg.logger.Info("Transaction confirmed",
		"txHash", tx.Hash().Hex(),
		"status", receipt.Status,
		"gasUsed", receipt.GasUsed,
		"blockNumber", receipt.BlockNumber,
		"blockHash", receipt.BlockHash.Hex())

	if receipt.Status == 0 {
		return fmt.Errorf("transaction failed: %s", tx.Hash().Hex())
	}
	return nil
}
