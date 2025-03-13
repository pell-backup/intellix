package taskgateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ChainConnection struct {
	ethClient          *ethclient.Client
	contractDataOracle *contractdataoracle.ContractDataOracle
}

type TaskGateway struct {
	service.BaseService

	server     *rpc.Server
	serverAddr string
	listener   net.Listener

	cfg *TaskGatewayCfg
	ctx context.Context

	logger     log.Logger
	privateKey *keystore.Key

	chainConnections map[uint64]*ChainConnection
	taskMap          sync.Map
	nonceMap         sync.Map
}

func NewTaskGateway(logger log.Logger, ctx context.Context, cfg *TaskGatewayCfg) (*TaskGateway, error) {
	logger = logger.With("comp", "gateway")
	chainConns := make(map[uint64]*ChainConnection)
	for chainID, chainCfg := range cfg.Chains {
		ethClient, err := ethclient.Dial(chainCfg.RPCURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to chain %d: %v", chainID, err)
		}

		contract, err := contractdataoracle.NewContractDataOracle(common.HexToAddress(chainCfg.ContractAddress), ethClient)
		if err != nil {
			return nil, fmt.Errorf("failed to create contract instance for chain %d: %v", chainID, err)
		}

		chainConns[chainID] = &ChainConnection{
			ethClient:          ethClient,
			contractDataOracle: contract,
		}
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
		server:           server,
		cfg:              cfg,
		ctx:              ctx,
		logger:           logger,
		privateKey:       key,
		serverAddr:       cfg.ServerAddr,
		chainConnections: chainConns,
		taskMap:          sync.Map{},
		nonceMap:         sync.Map{},
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
	for chainID, conn := range tg.chainConnections {
		conn.ethClient.Close()
		tg.logger.Info("Closed connection", "chainID", chainID)
	}
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
	chainConn, ok := tg.chainConnections[uint64(response.ChainID)]
	if !ok {
		return fmt.Errorf("no connection found for chain ID %d", response.ChainID)
	}

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
	task := contractdataoracle.IDataOracleTask{
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

	sign := contractdataoracle.IBLSSignatureVerifierNonSignerStakesAndSignature{
		NonSignerGroupBitmapIndices: response.ValidatedData.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:            convertNonSignersPubkeysG1(response.ValidatedData.NonSignersPubkeysG1),
		GroupApks:                   convertQuorumApks(response.ValidatedData.QuorumApksG1),
		ApkG2:                       convertApkG2(response.ValidatedData.SignersApkG2),
		Sigma:                       convertSigma(response.ValidatedData.SignersAggSigG1),
		GroupApkIndices:             response.ValidatedData.QuorumApkIndices,
		TotalStakeIndices:           response.ValidatedData.TotalStakeIndices,
		NonSignerStakeIndices:       response.ValidatedData.NonSignerStakeIndices,
	}

	taskResp := contractdataoracle.IDataOracleTaskResponse{
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

	if conf, ok := tg.cfg.Chains[uint64(response.ChainID)]; ok {
		if conf.GasLimit > 0 {
			authOpts.GasLimit = conf.GasLimit
		}
	}

	transaction, err := chainConn.contractDataOracle.ResponseToTask(authOpts, task, taskResp, sign)
	if err != nil {
		tg.logger.Error("Error assembling RequestPrice tx",
			"chainID", response.ChainID,
			"err", err,
			"task", fmt.Sprintf("%+v", task),
			"taskResp", fmt.Sprintf("%+v", taskResp),
			"sign", fmt.Sprintf("%+v", sign))
		return err
	}

	if transaction != nil {
		err = tg.queryTransaction(ctx, transaction, chainConn.ethClient)
		if err != nil {
			return fmt.Errorf("chain %d: %v", response.ChainID, err)
		}
	}

	return nil
}

func (tg *TaskGateway) queryTransaction(ctx context.Context, tx *types.Transaction, ethClient *ethclient.Client) error {
	tg.logger.Info("Transaction submitted", "txHash", tx.Hash().Hex())

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	receipt, err := bind.WaitMined(timeoutCtx, ethClient, tx)
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
