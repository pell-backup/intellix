package taskgateway

import (
	"context"
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

type TaskGateway struct {
	service.BaseService

	logger log.Logger
	client *ethclient.Client

	contractPriceOracle *priceOracle.ContractPriceOracle

	responseChan chan *types.MsgVoteFinalizedRequestPrice
	taskMap      sync.Map
}

func NewTaskGateway(logger log.Logger, ethEndpoint string) (*TaskGateway, error) {
	client, err := ethclient.Dial(ethEndpoint)
	if err != nil {
		return nil, err
	}

	tg := &TaskGateway{
		logger:       logger,
		client:       client,
		responseChan: make(chan *types.MsgVoteFinalizedRequestPrice, 100),
	}

	tg.BaseService = *service.NewBaseService(logger, "TaskGateway", tg)
	return tg, nil
}

func (tg *TaskGateway) OnStart() error {
	go tg.processResponses(context.Background())
	go tg.listenForVoteFinalizedRequestPrice()
	return nil
}

func (tg *TaskGateway) listenForVoteFinalizedRequestPrice() {
	for {
		select {
		case <-tg.Quit():
			return
		default:
			event, err := tg.SubscribeToVoteFinalizedRequestPrice()
			if err != nil {
				tg.logger.Error("Error subscribing to VoteFinalizedRequestPrice", "err", err)
				time.Sleep(2 * time.Second) // retry
				continue
			}
			tg.handleVoteFinalizedRequestPrice(event)
		}
	}
}

func (tg *TaskGateway) SubscribeToVoteFinalizedRequestPrice() (*types.MsgVoteFinalizedRequestPrice, error) {
	// listen chain MsgVoteFinalizedRequestPrice

	return nil, nil
}

func (tg *TaskGateway) handleVoteFinalizedRequestPrice(event *types.MsgVoteFinalizedRequestPrice) {
	tg.SubmitResponse(event)
}

func (tg *TaskGateway) OnStop() {
	close(tg.responseChan)
}

func (tg *TaskGateway) SubmitResponse(response *types.MsgVoteFinalizedRequestPrice) {
	tg.responseChan <- response
}

func (tg *TaskGateway) processResponses(ctx context.Context) {
	for response := range tg.responseChan {
		tg.handleResponse(ctx, response)
	}
}

func (tg *TaskGateway) handleResponse(ctx context.Context, response *types.MsgVoteFinalizedRequestPrice) {
	value, loaded := tg.taskMap.LoadOrStore(response.RequestId, response)
	if loaded {
		existingResponse := value.(*types.MsgVoteFinalizedRequestPrice)
		if tg.shouldReplaceResponse(existingResponse, response) {
			tg.taskMap.Store(response.RequestId, response)
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

func (tg *TaskGateway) serializeSubmitChainData(response *types.MsgVoteFinalizedRequestPrice) ([]byte, error) {
	// abi encode
	return nil, nil
}

func (tg *TaskGateway) submitToChain(ctx context.Context, response *types.MsgVoteFinalizedRequestPrice) {
	var data, err = tg.serializeSubmitChainData(response)
	if err != nil {
		tg.logger.Error("Error serializeSubmitChainData", "err", err)
		return
	}
	var nonce *big.Int
	var sender common.Address
	txOpts := &bind.TransactOpts{}

	feeTokenAddr, err := common.NewMixedcaseAddressFromString(response.FeeToken)
	if err != nil {
		tg.logger.Error("Error NewMixedcaseAddressFromString", "err", err)
		return
	}
	var payment = response.Payment.BigInt()
	var callbackFuncId = [4]byte(response.CallbackFunctionId)
	var quorumThresholdPercentage = response.QuorumThresholdPercentage
	var quorumNumbers = response.QuorumNumbers

	transaction, err := tg.contractPriceOracle.RequestPrice(txOpts, sender, payment, feeTokenAddr.Address(), txOpts.From, callbackFuncId, nonce, data, quorumThresholdPercentage, quorumNumbers)
	if err != nil {
		tg.logger.Error("Error assembling RequestPrice tx", "err", err)
		return
	}

	receipt, err := tg.sendTransaction(ctx, transaction)
	if err != nil {
		tg.logger.Error("Error sending transaction", "err", err)
		return
	}

	var receiptLog = receipt.Logs[1]
	newTaskCreatedEvent, err := tg.contractPriceOracle.ParseNewTaskCreated(*receiptLog)
	if err != nil {
		tg.logger.Error("Error parsing NewTaskCreated event", "err", err)
		return
	}

	_ = newTaskCreatedEvent

	return
}

func (tg *TaskGateway) sendTransaction(ctx context.Context, transaction *ethtypes.Transaction) (*ethtypes.Receipt, error) {
	// send transaction to chain and get log
	err := tg.client.SendTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return tg.client.TransactionReceipt(ctx, transaction.Hash())
}
