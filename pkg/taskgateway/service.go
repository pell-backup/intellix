package taskgateway

import (
	"context"
	"fmt"
	tmclient "github.com/cometbft/cometbft/rpc/client/http"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	priceOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/PriceOracle"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"intellix/x/price/types"

	"sync"
	"time"
)

var tmClientQuery = fmt.Sprintf("tm.event='%s'", types.EventTypeVoteFinalizedRequestPrice)

type TaskGateway struct {
	service.BaseService

	cfg *TaskGatewayCfg
	ctx context.Context

	logger    log.Logger
	ethClient *ethclient.Client
	tmClient  *tmclient.HTTP
	cdc       codec.Codec

	contractPriceOracle *priceOracle.ContractPriceOracle
	responseChan        chan *types.MsgVoteFinalizedRequestPrice
	taskMap             sync.Map
	nonceMap            sync.Map
}

func NewTaskGateway(logger log.Logger, ctx context.Context, cfg *TaskGatewayCfg) (*TaskGateway, error) {
	ethClient, err := ethclient.Dial(cfg.EthEndpoint)
	if err != nil {
		return nil, err
	}

	tmClient, err := tmclient.New(cfg.BftNetworkRemote, "/websocket")
	if err != nil {
		return nil, err
	}

	contract, err := priceOracle.NewContractPriceOracle(common.HexToAddress(cfg.ContractAddress), ethClient)
	if err != nil {
		return nil, err
	}

	tg := &TaskGateway{
		cfg:                 cfg,
		ctx:                 ctx,
		logger:              logger,
		ethClient:           ethClient,
		tmClient:            tmClient,
		cdc:                 newCodec(),
		contractPriceOracle: contract,
		responseChan:        make(chan *types.MsgVoteFinalizedRequestPrice, 100),
		taskMap:             sync.Map{},
		nonceMap:            sync.Map{},
	}

	tg.BaseService = *service.NewBaseService(logger, "TaskGateway", tg)
	return tg, nil
}

func newCodec() codec.Codec {
	registry := sdktypes.NewInterfaceRegistry()
	registry.RegisterImplementations((*sdk.Msg)(nil), &types.MsgVoteFinalizedRequestPrice{})
	return codec.NewProtoCodec(registry)
}

func (tg *TaskGateway) OnStart() error {
	err := tg.tmClient.Start()
	if err != nil {
		return err
	}
	go tg.processResponses(tg.ctx)
	go tg.listenForVoteFinalizedRequestPrice()
	return nil
}

func (tg *TaskGateway) listenForVoteFinalizedRequestPrice() {
	eventCh, err := tg.tmClient.Subscribe(tg.ctx, "task-gateway", tmClientQuery)
	if err != nil {
		tg.logger.Error("Error subscribing to VoteFinalizedRequestPrice", "err", err)
		return
	}

	for {
		select {
		case <-tg.Quit():
			return
		case event := <-eventCh:
			msg, err := tg.convertEventMsgVoteFinalizedRequestPrice(event.Data)
			if err != nil {
				tg.logger.Error("Error subscribing to VoteFinalizedRequestPrice", "err", err)
				time.Sleep(2 * time.Second) // retry
				continue
			}
			tg.responseChan <- msg
		}
	}
}

func (tg *TaskGateway) convertEventMsgVoteFinalizedRequestPrice(eventData tmtypes.TMEventData) (*types.MsgVoteFinalizedRequestPrice, error) {
	tx, ok := eventData.(tmtypes.EventDataTx)
	if !ok {
		return nil, fmt.Errorf("expected EventDataTx, got %T", eventData)
	}

	// decode tx
	var out = &types.MsgVoteFinalizedRequestPrice{}
	err := tg.cdc.Unmarshal(tx.Tx, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (tg *TaskGateway) OnStop() {
	_ = tg.tmClient.Stop()
	tg.ethClient.Close()
	close(tg.responseChan)
}

func (tg *TaskGateway) processResponses(ctx context.Context) {
	for response := range tg.responseChan {
		tg.handleResponse(ctx, response)
	}
}

func (tg *TaskGateway) handleResponse(ctx context.Context, response *types.MsgVoteFinalizedRequestPrice) {
	value, loaded := tg.taskMap.LoadOrStore(response.TaskRaw.TaskIndex, response)
	if loaded {
		existingResponse := value.(*types.MsgVoteFinalizedRequestPrice)
		if tg.shouldReplaceResponse(ctx, existingResponse, response) {
			tg.taskMap.Store(response.TaskRaw.RequestId, response)
			go tg.submitToChain(ctx, response)
		}
	} else {
		go tg.submitToChain(ctx, response)
	}
}

func (tg *TaskGateway) shouldReplaceResponse(ctx context.Context, existing, new *types.MsgVoteFinalizedRequestPrice) bool {
	// TODO: add security threshold comparison
	return false
}

func (tg *TaskGateway) submitToChain(ctx context.Context, response *types.MsgVoteFinalizedRequestPrice) {
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
		return
	}
	cbAddr, err := convertAddressToString(response.TaskRaw.CallbackAddress)
	if err != nil {
		tg.logger.Error("Error converting callback address", "err", err)
		return
	}

	transaction, err := tg.contractPriceOracle.UpdatePrice(txOpts, priceOracle.IPriceOracleTask{
		RequestId:                 [32]byte(response.TaskRaw.RequestId),
		FeeToken:                  *feeTokenAddr,
		Payment:                   response.TaskRaw.Payment.BigInt(),
		RequestData:               response.TaskRaw.RequestData,
		CallbackAddress:           *cbAddr,
		CallbackFunctionId:        [4]byte(response.TaskRaw.CallbackFunctionId),
		TaskCreatedBlock:          response.TaskRaw.TaskCreatedBlock,
		QuorumNumbers:             response.TaskRaw.QuorumNumbers,
		QuorumThresholdPercentage: response.TaskRaw.QuorumThresholdPercentage,
	}, priceOracle.IPriceOracleTaskResponse{
		ReferenceTaskIndex: response.PriceFeedResponse.ReferenceTaskIndex,
		Price:              response.PriceFeedResponse.Price.BigInt(),
	}, priceOracle.IBLSSignatureCheckerNonSignerStakesAndSignature{
		NonSignerQuorumBitmapIndices: response.ValidatedData.NonSignerQuorumBitmapIndices,
		NonSignerPubkeys:             convertPbToBN254G1PointList(response.ValidatedData.NonSignersPubkeysG1),
		QuorumApks:                   convertPbToBN254G1PointList(response.ValidatedData.QuorumApksG1),
		ApkG2:                        *convertPbToBN254G2Point(response.ValidatedData.SignersApkG2),
		Sigma:                        *convertPbToBN254G1Point(response.ValidatedData.SignersAggSigG1),
		QuorumApkIndices:             response.ValidatedData.QuorumApkIndices,
		TotalStakeIndices:            response.ValidatedData.TotalStakeIndices,
		NonSignerStakeIndices:        convertUInt32ListToSlice(response.ValidatedData.NonSignerStakeIndices),
	})
	if err != nil {
		tg.logger.Error("Error assembling RequestPrice tx", "err", err)
		return
	}

	receipt, err := tg.sendTransaction(ctx, transaction)
	if err != nil {
		tg.logger.Error("Error sending transaction", "err", err)
		return
	}
	_ = receipt
}

func (tg *TaskGateway) sendTransaction(ctx context.Context, transaction *ethtypes.Transaction) (*ethtypes.Receipt, error) {
	// send transaction to chain and get log
	err := tg.ethClient.SendTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return tg.ethClient.TransactionReceipt(ctx, transaction.Hash())
}
