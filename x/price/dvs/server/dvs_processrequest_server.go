package server

import (
	"bytes"
	"context"
	"cosmossdk.io/math"
	"fmt"
	cmttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"intellix/x/price/dvs/types"
	"math/big"
	"time"

	"sort"
)

type DvsProcessRequestServer struct {
	Server
}

// NewDvsProcessRequestServer returns an implementation of the DvsProcessRequestServer interface
// for the provided Server.
func NewDvsProcessRequestServer(server Server) types.DvsProcessRequestServer {
	return &DvsProcessRequestServer{
		Server: server,
	}
}

func (d *DvsProcessRequestServer) ProcessRequestPriceFeed(ctx context.Context, request *types.ProcessRequestPriceFeedIn) (*types.ProcessRequestPriceFeedOut, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// just for testing
	//{
	//	return &types.ProcessRequestPriceFeedOut{
	//		RequestId:     request.Task.RequestId,
	//		TaskIndex:     request.Task.TaskIndex,
	//		Price:         math.LegacyNewDec(1),
	//		Timestamp:     time.Now().Unix(),
	//		SourceCount:   0,
	//		OperatorCount: 0,
	//		BlockHeight:   0,
	//		BlockRange:    nil,
	//	}, nil
	//}

	// fetch raw price from chain
	rawPrices, err := fetchRawPrices(sdkCtx, d.Logger(), request.PriceFeed.BaseSymbol, request.PriceFeed.QuoteSymbol)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch raw prices: %w", err)
	}

	// sign data and broadcast VoteRequestPriceFeed
	err = d.broadcastVoteRequestPriceFeed(sdkCtx, request, request.GetPriceFeed(), rawPrices)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast VoteRequestPriceFeed: %w", err)
	}

	// listen and collect [N-N+M] block
	priceFeedTxs, err := d.collectVoteRequestPriceFeedTxs(sdkCtx, request.Task.RequestId)
	if err != nil {
		return nil, fmt.Errorf("failed to collect VoteRequestPriceFeed transactions: %w", err)
	}

	// aggregate [N-N+M] block prices
	return d.aggregatePrices(sdkCtx, request.Task.TaskIndex, request.Task.RequestId, priceFeedTxs)
}

func (d *DvsProcessRequestServer) broadcastVoteRequestPriceFeed(ctx sdk.Context, task *types.ProcessRequestPriceFeedIn, priceFeed *types.PriceFeedParam, rawPrices map[string]*big.Int) error {
	for dataSource, price := range rawPrices {
		msg := types.MsgVoteRequestPriceFeed{
			TaskIndex:   task.Task.TaskIndex,
			OperatorId:  d.Server.GetOperatorAddress(ctx),
			RequestId:   task.Task.RequestId,
			BaseSymbol:  priceFeed.BaseSymbol,
			QuoteSymbol: priceFeed.QuoteSymbol,
			Price:       math.LegacyNewDecFromBigInt(price),
			Timestamp:   time.Now().Unix(),
			BlockHeight: uint64(ctx.BlockHeight()),
		}

		if err := d.Server.SignAndBroadcastTx(ctx, &msg); err != nil {
			return fmt.Errorf("failed to broadcast VoteRequestPriceFeed for data source %s: %w", dataSource, err)
		}
	}
	return nil
}

func (d *DvsProcessRequestServer) collectVoteRequestPriceFeedTxs(ctx sdk.Context, requestID []byte) ([]types.MsgVoteRequestPriceFeed, error) {
	var priceFeedTxs []types.MsgVoteRequestPriceFeed
	firstTxBlock := int64(0)

	for {
		block, err := d.Server.GetLatestBlock(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get latest block: %w", err)
		}

		newTxs := d.processBlockTxs(ctx, block, requestID)
		priceFeedTxs = append(priceFeedTxs, newTxs...)

		if firstTxBlock == 0 && len(newTxs) > 0 {
			firstTxBlock = block.Header.Height
		}

		// check N-N+M
		if d.shouldStopCollecting(ctx, firstTxBlock, block.Header.Height) {
			break
		}

		// wait for next block
		ctx = ctx.WithBlockHeight(block.Header.Height + 1)
	}

	return priceFeedTxs, nil
}

func (d *DvsProcessRequestServer) processBlockTxs(ctx sdk.Context, block *cmttypes.Block, requestID []byte) []types.MsgVoteRequestPriceFeed {
	var priceFeedTxs []types.MsgVoteRequestPriceFeed

	for _, tx := range block.Data.Txs {
		// only collect VoteRequestPriceFeed && current requestId data
		if msg, ok := d.isVoteRequestPriceFeedTx(tx); ok && bytes.Equal(msg.RequestId, requestID) {
			priceFeedTxs = append(priceFeedTxs, types.MsgVoteRequestPriceFeed{
				TaskIndex:   msg.TaskIndex,
				OperatorId:  msg.OperatorId,
				RequestId:   msg.RequestId,
				BaseSymbol:  msg.BaseSymbol,
				QuoteSymbol: msg.QuoteSymbol,
				Price:       msg.Price,
				Timestamp:   msg.Timestamp,
				BlockHeight: uint64(block.Header.Height),
			})
		}
	}

	return priceFeedTxs
}

func (d *DvsProcessRequestServer) isVoteRequestPriceFeedTx(tx cmttypes.Tx) (*types.MsgVoteRequestPriceFeed, bool) {
	// check tx is VoteRequestPriceFeed
	decoder := d.clientCtx.TxConfig.TxDecoder()
	data, err := decoder(tx)
	if err != nil {
		d.Logger().Error("TxDecoder decode tx error: " + err.Error())
		return nil, false
	}

	msgs := data.GetMsgs()
	if len(msgs) == 0 {
		return nil, false
	}

	msg := msgs[0]
	voteMsg, ok := msg.(*types.MsgVoteRequestPriceFeed)
	if !ok {
		return nil, false
	}

	return voteMsg, true
}

func (d *DvsProcessRequestServer) shouldStopCollecting(ctx sdk.Context, firstTxBlock, currentBlock int64) bool {
	if firstTxBlock == 0 {
		return false
	}
	return currentBlock >= firstTxBlock+d.waitBlockCount
}

func (d *DvsProcessRequestServer) aggregatePrices(ctx sdk.Context, taskIndex uint32, requestID []byte, priceFeedTxs []types.MsgVoteRequestPriceFeed) (*types.ProcessRequestPriceFeedOut, error) {
	operatorPrices := make(map[string][]math.LegacyDec)

	// calc avg the prices of different data sources within each Operator
	for _, tx := range priceFeedTxs {
		operatorID := tx.OperatorId
		operatorPrices[operatorID] = append(operatorPrices[operatorID], tx.Price)
	}

	var averagePrices []math.LegacyDec
	for _, prices := range operatorPrices {
		sum := math.LegacyZeroDec()
		for _, price := range prices {
			sum = sum.Add(price)
		}

		average := sum.Quo(math.LegacyNewDec(int64(len(prices))))
		averagePrices = append(averagePrices, average)
	}

	// median of the average prices of different Operators
	sort.Slice(averagePrices, func(i, j int) bool {
		return averagePrices[i].LT(averagePrices[j])
	})

	medianPrice := averagePrices[len(averagePrices)/2]

	blockRange := &types.BlockRange{}
	if len(priceFeedTxs) > 0 {
		blockRange.Start = priceFeedTxs[0].BlockHeight
		blockRange.End = priceFeedTxs[len(priceFeedTxs)-1].BlockHeight
	}

	return &types.ProcessRequestPriceFeedOut{
		RequestId:     requestID,
		TaskIndex:     taskIndex,
		Price:         medianPrice,
		Timestamp:     ctx.BlockTime().Unix(),
		SourceCount:   int32(len(priceFeedTxs)),
		OperatorCount: int32(len(operatorPrices)),
		BlockHeight:   uint64(ctx.BlockHeight()),
		BlockRange:    blockRange,
	}, nil
}
