package taskgateway

import (
	"context"
	"fmt"
	tmclient "github.com/cometbft/cometbft/rpc/client/http"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdktypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"intellix/x/price/dvs/types"
	"math/big"

	priceOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/PriceOracle"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/libs/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"sync"
	"time"
)

const tmClientQuery = "tm.event='Tx' AND message.action='/intellix.price.dvs.MsgVoteFinalizedRequestPrice'"

type TaskGateway struct {
	service.BaseService

	cfg *TaskGatewayCfg
	ctx context.Context

	logger    log.Logger
	ethClient *ethclient.Client
	tmClient  *tmclient.HTTP
	clientCtx *client.Context // just for decode tx

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

	tmClient, err := tmclient.New(cfg.CosmosNetworkUrl, "/websocket")
	if err != nil {
		return nil, err
	}

	contract, err := priceOracle.NewContractPriceOracle(common.HexToAddress(cfg.ContractAddress), ethClient)
	if err != nil {
		return nil, err
	}

	clientCtx, err := registerClientCtx()
	if err != nil {
		return nil, err
	}

	tg := &TaskGateway{
		cfg:                 cfg,
		ctx:                 ctx,
		logger:              logger,
		ethClient:           ethClient,
		tmClient:            tmClient,
		clientCtx:           clientCtx,
		contractPriceOracle: contract,
		responseChan:        make(chan *types.MsgVoteFinalizedRequestPrice, 100),
		taskMap:             sync.Map{},
		nonceMap:            sync.Map{},
	}

	tg.BaseService = *service.NewBaseService(logger, "TaskGateway", tg)
	return tg, nil
}

func registerClientCtx() (*client.Context, error) {
	// register client context
	registry := sdktypes.NewInterfaceRegistry()
	registry.RegisterImplementations((*sdk.Msg)(nil), &types.MsgVoteFinalizedRequestPrice{})

	clientCtx := client.Context{}.WithCodec(
		codec.NewProtoCodec(registry),
	)

	return &clientCtx, nil
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
		default:
			// pass
		}
	}
}

func (tg *TaskGateway) convertEventMsgVoteFinalizedRequestPrice(eventData tmtypes.TMEventData) (*types.MsgVoteFinalizedRequestPrice, error) {
	tx, ok := eventData.(tmtypes.EventDataTx)
	if !ok {
		return nil, fmt.Errorf("expected EventDataTx, got %T", eventData)
	}

	// decode tx
	decoder := tg.clientCtx.TxConfig.TxDecoder()
	data, err := decoder(tx.GetTx())
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx: %w", err)
	}
	if len(data.GetMsgs()) != 1 {
		return nil, fmt.Errorf("expected 1 message, got %d", len(data.GetMsgs()))
	}
	out, ok := data.GetMsgs()[0].(*types.MsgVoteFinalizedRequestPrice)
	if !ok {
		return nil, fmt.Errorf("expected MsgVoteFinalizedRequestPrice, got %T", data.GetMsgs()[0])
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
	value, loaded := tg.taskMap.LoadOrStore(response.Task.TaskIndex, response)
	if loaded {
		existingResponse := value.(*types.MsgVoteFinalizedRequestPrice)
		if tg.shouldReplaceResponse(existingResponse, response) {
			tg.taskMap.Store(response.Task.RequestId, response)
			go tg.submitToChain(ctx, response)
		}
	} else {
		go tg.submitToChain(ctx, response)
	}
}

func (tg *TaskGateway) shouldReplaceResponse(existing, new *types.MsgVoteFinalizedRequestPrice) bool {
	// TODO: add security threshold comparison
	return false
}

func (tg *TaskGateway) getNonce(ctx context.Context, address string) (*big.Int, error) {
	if nonce, ok := tg.nonceMap.Load(address); ok {
		nonce.(*big.Int).Add(nonce.(*big.Int), big.NewInt(1))
		return nonce.(*big.Int), nil
	}

	nonce, err := tg.ethClient.PendingNonceAt(ctx, common.HexToAddress(address))
	if err != nil {
		return nil, fmt.Errorf("getNonce err: %w", err)
	}
	outNonce := big.NewInt(int64(nonce + 1))
	tg.nonceMap.Store(address, outNonce)

	return outNonce, nil
}

func (tg *TaskGateway) submitToChain(ctx context.Context, response *types.MsgVoteFinalizedRequestPrice) {
	var sender = common.HexToAddress(tg.cfg.ContractFromAddress)
	txOpts := &bind.TransactOpts{
		From: sender,
		Signer: func(common.Address, *ethtypes.Transaction) (*ethtypes.Transaction, error) {
			return nil, nil
		},
	}

	feeTokenAddr, err := convertAddressToString(response.Task.FeeToken)
	if err != nil {
		tg.logger.Error("Error converting fee token address", "err", err)
		return
	}
	cbAddr, err := convertAddressToString(response.Task.CallbackAddress)
	if err != nil {
		tg.logger.Error("Error converting callback address", "err", err)
		return
	}

	transaction, err := tg.contractPriceOracle.UpdatePrice(txOpts, priceOracle.IPriceOracleTask{
		RequestId:                 [32]byte(response.Task.RequestId),
		FeeToken:                  *feeTokenAddr,
		Payment:                   response.Task.Payment.BigInt(),
		RequestData:               response.Task.RequestData,
		CallbackAddress:           *cbAddr,
		CallbackFunctionId:        [4]byte(response.Task.CallbackFunctionId),
		TaskCreatedBlock:          response.Task.TaskCreatedBlock,
		QuorumNumbers:             response.Task.QuorumNumbers,
		QuorumThresholdPercentage: response.Task.QuorumThresholdPercentage,
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

	return
}

func (tg *TaskGateway) sendTransaction(ctx context.Context, transaction *ethtypes.Transaction) (*ethtypes.Receipt, error) {
	// send transaction to chain and get log
	err := tg.ethClient.SendTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return tg.ethClient.TransactionReceipt(ctx, transaction.Hash())
}
