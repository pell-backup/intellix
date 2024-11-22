package taskgateway

import (
	"context"
	"fmt"
	dvslog "github.com/0xPellNetwork/pelldvs/libs/log"
	"net"
	"net/rpc"

	dataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"intellix/x/price/types"

	"sync"
)

type TaskGateway struct {
	service.BaseService

	server     *rpc.Server
	serverAddr string
	listener   net.Listener

	cfg *TaskGatewayCfg
	ctx context.Context

	logger    dvslog.Logger
	ethClient *ethclient.Client

	contractDataOracle *dataOracle.ContractDataOracle
	taskMap            sync.Map
	nonceMap           sync.Map
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

	server := rpc.NewServer()
	tg := &TaskGateway{
		server:             server,
		cfg:                cfg,
		ctx:                ctx,
		logger:             logger,
		ethClient:          ethClient,
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

func (tg *TaskGateway) RespondToTask(req *types.MsgVoteFinalizedRequestPrice, resp *RespondToTaskResponse) error {
	err := tg.handleResponse(context.Background(), req)
	if err != nil {
		resp.Error = err.Error()
		return err
	}
	resp.Error = ""
	return nil
}

func (tg *TaskGateway) handleResponse(ctx context.Context, response *types.MsgVoteFinalizedRequestPrice) error {
	value, loaded := tg.taskMap.LoadOrStore(response.TaskRaw.TaskIndex, response)
	if loaded {
		existingResponse := value.(*types.MsgVoteFinalizedRequestPrice)
		if tg.shouldReplaceResponse(ctx, existingResponse, response) {
			tg.taskMap.Store(response.TaskRaw.RequestId, response)
			return tg.submitToChain(ctx, response)
		}
	} else {
		return tg.submitToChain(ctx, response)
	}
	return nil
}

func (tg *TaskGateway) shouldReplaceResponse(ctx context.Context, existing, new *types.MsgVoteFinalizedRequestPrice) bool {
	// TODO: add security threshold comparison
	return false
}

func (tg *TaskGateway) submitToChain(ctx context.Context, response *types.MsgVoteFinalizedRequestPrice) error {
	var sender = common.HexToAddress(tg.cfg.SenderAddress)
	txOpts := &bind.TransactOpts{
		From: sender,
		Signer: func(common.Address, *ethtypes.Transaction) (*ethtypes.Transaction, error) {
			return nil, nil
		},
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

	transaction, err := tg.contractDataOracle.ResponseToTask(txOpts, dataOracle.IDataOracleTask{
		RequestId:                 [32]byte(response.TaskRaw.RequestId),
		FeeToken:                  *feeTokenAddr,
		Payment:                   response.TaskRaw.Payment.BigInt(),
		RequestData:               response.TaskRaw.RequestData,
		CallbackAddress:           *cbAddr,
		CallbackFunctionId:        [4]byte(response.TaskRaw.CallbackFunctionId),
		TaskCreatedBlock:          response.TaskRaw.TaskCreatedBlock,
		QuorumNumbers:             response.TaskRaw.QuorumNumbers,
		QuorumThresholdPercentage: response.TaskRaw.QuorumThresholdPercentage,
	}, dataOracle.IDataOracleTaskResponse{
		ReferenceTaskIndex: response.PriceFeedResponse.ReferenceTaskIndex,
		Price:              response.PriceFeedResponse.Price.BigInt(),
	}, dataOracle.IBLSSignatureCheckerNonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: response.ValidatedData.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:             convertNonSignersPubkeysG1(response.ValidatedData.NonSignersPubkeysG1),
		QuorumApks:                   convertQuorumApks(response.ValidatedData.QuorumApksG1),
		ApkG2:                        *convertApkG2(response.ValidatedData.SignersApkG2),
		Sigma:                        *convertSigma(response.ValidatedData.SignersAggSigG1),
		QuorumApkIndices:             response.ValidatedData.QuorumApkIndices,
		TotalStakeIndices:            response.ValidatedData.TotalStakeIndices,
		NonSignerStakeIndices:        convertNonSignerStakeIndices(response.ValidatedData.NonSignerStakeIndices),
	})
	if err != nil {
		tg.logger.Error("Error assembling RequestPrice tx", "err", err)
		return err
	}

	receipt, err := tg.sendTransaction(ctx, transaction)
	if err != nil {
		tg.logger.Error("Error sending transaction", "err", err)
		return err
	}
	_ = receipt
	return nil
}

func (tg *TaskGateway) sendTransaction(ctx context.Context, transaction *ethtypes.Transaction) (*ethtypes.Receipt, error) {
	// send transaction to chain and get log
	err := tg.ethClient.SendTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return tg.ethClient.TransactionReceipt(ctx, transaction.Hash())
}
