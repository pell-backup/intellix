package server

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"time"

	"cosmossdk.io/math"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	"intellix/dvs/price/types"
	pricetypes "intellix/x/price/types"
)

func (s Server) RequestPriceFeed(ctx context.Context, request *types.RequestPriceFeedIn) (*types.RequestPriceFeedOut, error) {
	pkgContext := sdktypes.UnwrapContext(ctx)
	s.logger.Info("ProcessRequestPriceFeed", "PriceFeedParam", fmt.Sprintf("%+v", request.PriceFeed))

	// fetch raw price from chain
	rawPrices, err := fetchRawPrices(pkgContext, s.Logger(), request.PriceFeed.BaseSymbol, request.PriceFeed.QuoteSymbol,
		ToPriceTickConverterByDataSource(s.tickConverterConfig), s.apiKey, s.symbolList)
	if err != nil {
		s.logger.Error("ProcessRequestPriceFeed fetchRawPrices error: " + err.Error())
		return nil, fmt.Errorf("failed to fetch raw prices: %w", err)
	}

	// sign data and broadcast VoteRequestPriceFeed
	err = s.broadcastVoteRequestPriceFeed(pkgContext, request, request.GetPriceFeed(), rawPrices)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast VoteRequestPriceFeed: %w", err)
	}

	// listen and collect [N-N+M] block
	priceFeedTxs, err := s.collectVoteRequestPriceFeedTxs(pkgContext, request.Task.RequestId)
	if err != nil {
		return nil, fmt.Errorf("failed to collect VoteRequestPriceFeed transactions: %w", err)
	}

	// aggregate [N-N+M] block prices
	return s.aggregatePrices(pkgContext, request.Task.TaskIndex, request.Task.RequestId, priceFeedTxs)
}

func (s Server) broadcastVoteRequestPriceFeed(ctx sdktypes.Context, task *types.RequestPriceFeedIn, priceFeed *types.PriceFeedParam, rawPrices map[string]math.LegacyDec) error {
	s.Logger().Info("broadcastVoteRequestPriceFeed",
		"rawPrices", fmt.Sprintf("%+v", rawPrices),
		"task", fmt.Sprintf("%+v", task),
		"priceFeed", fmt.Sprintf("%+v", priceFeed),
	)

	var prices []*pricetypes.VoteRequestPriceFeed
	for dataSource, price := range rawPrices {
		prices = append(prices, &pricetypes.VoteRequestPriceFeed{
			Source: dataSource,
			Price:  price,
		})
	}
	addr, err := s.SenderAddress()
	if err != nil {
		return err
	}

	msg := pricetypes.MsgVoteRequestPriceFeed{
		Sender:      addr.String(),
		TaskIndex:   task.Task.TaskIndex,
		OperatorId:  s.GetOperatorAddress(ctx),
		RequestId:   task.Task.RequestId,
		BaseSymbol:  priceFeed.BaseSymbol,
		QuoteSymbol: priceFeed.QuoteSymbol,
		Price:       prices,
		Timestamp:   time.Now().Unix(),
		BlockHeight: uint64(ctx.Height()),
	}
	if err := s.SignAndBroadcastTx(ctx, &msg); err != nil {
		s.Logger().Error("broadcastVoteRequestPriceFeed SignAndBroadcastTx error: " + err.Error())
		return fmt.Errorf("failed to broadcast VoteRequestPriceFeed for data error: %w", err)
	}

	return nil
}

func (s Server) collectVoteRequestPriceFeedTxs(ctx sdktypes.Context, requestID []byte) ([]pricetypes.MsgVoteRequestPriceFeed, error) {
	var priceFeedTxs []pricetypes.MsgVoteRequestPriceFeed
	firstTxBlock := int64(0)

	for {
		block, err := s.GetLatestBlock(ctx)
		if err != nil {
			s.logger.Error("collectVoteRequestPriceFeed GetLatestBlock error: " + err.Error())
			return nil, fmt.Errorf("failed to get latest block: %w", err)
		}

		newTxs := s.processBlockTxs(ctx, block, requestID)
		priceFeedTxs = append(priceFeedTxs, newTxs...)

		if firstTxBlock == 0 && len(newTxs) > 0 {
			firstTxBlock = block.Header.Height
		}

		// check N-N+M
		if s.shouldStopCollecting(ctx, firstTxBlock, block.Header.Height) {
			s.logger.Info("collectVoteRequestPriceFeed stop collecting", "block_height", block.Header.Height)
			break
		}

		// wait for next block
		time.Sleep(time.Millisecond * 10)
	}

	return priceFeedTxs, nil
}

func (s Server) processBlockTxs(ctx context.Context, block *cmttypes.Block, requestID []byte) []pricetypes.MsgVoteRequestPriceFeed {
	var priceFeedTxs []pricetypes.MsgVoteRequestPriceFeed
	//s.logger.Info("Processing block", "height", block.Header.Height, "tx_count", len(block.Data.Txs))

	for _, tx := range block.Data.Txs {
		// only collect VoteRequestPriceFeed && current requestId data
		//s.logger.Info("Processing transaction", "tx_hash", tx.Hash())
		if msg, ok := s.isVoteRequestPriceFeedTx(tx); ok && bytes.Equal(msg.RequestId, requestID) {
			priceFeedTxs = append(priceFeedTxs, pricetypes.MsgVoteRequestPriceFeed{
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

func (s Server) isVoteRequestPriceFeedTx(tx cmttypes.Tx) (*pricetypes.MsgVoteRequestPriceFeed, bool) {
	// check tx is VoteRequestPriceFeed
	decoder := s.clientCtx.TxConfig.TxDecoder()
	data, err := decoder(tx)
	if err != nil {
		s.logger.Error("TxDecoder decode tx error: " + err.Error())
		return nil, false
	}

	msgs := data.GetMsgs()
	if len(msgs) == 0 {
		return nil, false
	}

	msg := msgs[0]

	// Try to handle authz message
	if authzMsg, ok := msg.(*authz.MsgExec); ok {
		// Get the inner messages from authz
		innerMsgs, err := authzMsg.GetMessages()
		if err != nil {
			s.logger.Error("Failed to get inner messages from authz", "error", err)
			return nil, false
		}
		if len(innerMsgs) == 0 {
			return nil, false
		}
		// Use the first inner message
		msg = innerMsgs[0]
	}

	voteMsg, ok := msg.(*pricetypes.MsgVoteRequestPriceFeed)
	if !ok {
		//s.logger.Error("msg is not MsgVoteRequestPriceFeed", "msg", fmt.Sprintf("%+v", msg))
		return nil, false
	}

	return voteMsg, true
}

func (s Server) shouldStopCollecting(ctx sdktypes.Context, firstTxBlock, currentBlock int64) bool {
	if firstTxBlock == 0 {
		return false
	}
	return currentBlock >= firstTxBlock+s.waitBlockCount
}

func (s Server) aggregatePrices(ctx sdktypes.Context, taskIndex uint32, requestID []byte, priceFeedTxs []pricetypes.MsgVoteRequestPriceFeed) (*types.RequestPriceFeedOut, error) {
	operatorPrices := make(map[string][]math.LegacyDec)

	// calc avg the prices of different data sources within each Operator
	for _, tx := range priceFeedTxs {
		operatorID := tx.OperatorId
		for _, price := range tx.Price {
			operatorPrices[operatorID] = append(operatorPrices[operatorID], price.Price)
		}
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

	return &types.RequestPriceFeedOut{
		RequestId:     requestID,
		TaskIndex:     taskIndex,
		Price:         medianPrice,
		Timestamp:     time.Now().Unix(),
		SourceCount:   int32(len(priceFeedTxs)),
		OperatorCount: int32(len(operatorPrices)),
		BlockHeight:   uint64(ctx.Height()),
		BlockRange:    blockRange,
	}, nil
}
