package server

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"cosmossdk.io/math"
<<<<<<< HEAD:dvs/price/server/request_server_price_feed.go
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	"github.com/cosmos/gogoproto/proto"

	"intellix/dvs/price/types"
	"intellix/pkg/tx_listener"
	sdktypes "intellix/sdk/types"
	"intellix/sdk/utils"
	pricetypes "intellix/x/price/types"
)

func (d *RequestServer) RequestPriceFeed(ctx context.Context, request *types.RequestPriceFeedIn) (*types.RequestPriceFeedOut, error) {
	pkgContext := sdktypes.UnwrapContext(ctx)
	d.logger.Info("ProcessRequestPriceFeed", "PriceFeedParam", fmt.Sprintf("%+v", request.PriceFeed))

	// fetch raw price from chain
	rawPrices, err := fetchRawPrices(pkgContext, d.Logger(), request.PriceFeed.BaseSymbol, request.PriceFeed.QuoteSymbol, ToPriceTickConverterByDataSource(d.tickConverterConfig))
	if err != nil {
		d.logger.Error("ProcessRequestPriceFeed fetchRawPrices error: " + err.Error())
=======
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
	rawPrices, err := fetchRawPrices(pkgContext, s.Logger(), request.PriceFeed.BaseSymbol, request.PriceFeed.QuoteSymbol, ToPriceTickConverterByDataSource(s.tickConverterConfig))
	if err != nil {
		s.logger.Error("ProcessRequestPriceFeed fetchRawPrices error: " + err.Error())
>>>>>>> main:dvs/price/server/msg_server_request_price_feed.go
		return nil, fmt.Errorf("failed to fetch raw prices: %w", err)
	}

	// sign data and broadcast VoteRequestPriceFeed
<<<<<<< HEAD:dvs/price/server/request_server_price_feed.go
	if err := d.broadcastVoteRequestPriceFeed(pkgContext, request, request.GetPriceFeed(), rawPrices); err != nil {
		d.logger.Error("ProcessRequestPriceFeed broadcastVoteRequestPriceFeed error: " + err.Error())
=======
	err = s.broadcastVoteRequestPriceFeed(pkgContext, request, request.GetPriceFeed(), rawPrices)
	if err != nil {
>>>>>>> main:dvs/price/server/msg_server_request_price_feed.go
		return nil, fmt.Errorf("failed to broadcast VoteRequestPriceFeed: %w", err)
	}

	// listen and collect [N-N+M] block
<<<<<<< HEAD:dvs/price/server/request_server_price_feed.go
	priceFeedTxs, err := d.startCollectEnoughPriceFeedTxs(pkgContext, request.Task.RequestId)
=======
	priceFeedTxs, err := s.collectVoteRequestPriceFeedTxs(pkgContext, request.Task.RequestId)
>>>>>>> main:dvs/price/server/msg_server_request_price_feed.go
	if err != nil {
		d.logger.Error("ProcessRequestPriceFeed collectVoteRequestPriceFeedTxs error: " + err.Error())
		return nil, fmt.Errorf("failed to collect VoteRequestPriceFeed transactions: %w", err)
	}
	if len(priceFeedTxs) == 0 {
		d.logger.Error("ProcessRequestPriceFeed collectVoteRequestPriceFeedTxs error: len(priceFeedTxs) == 0")
		return nil, fmt.Errorf("failed to collect VoteRequestPriceFeed transactions, len(priceFeedTxs) == 0")
	}
	d.logger.Info("ProcessRequestPriceFeed collectVoteRequestPriceFeedTxs success", "lengeth", len(priceFeedTxs), "priceFeedTxs", fmt.Sprintf("%+v", priceFeedTxs))

	// aggregate [N-N+M] block prices
<<<<<<< HEAD:dvs/price/server/request_server_price_feed.go
	return d.aggregatePrices(pkgContext, request.Task.TaskIndex, request.Task.RequestId, priceFeedTxs)
}

func (d *RequestServer) getMsgBytes(msg *pricetypes.MsgVoteRequestPriceFeed) []byte {
	m := &pricetypes.MsgVoteRequestPriceFeed{
		TaskIndex:   msg.TaskIndex,
		OperatorId:  msg.OperatorId,
		RequestId:   msg.RequestId,
		BaseSymbol:  msg.BaseSymbol,
		QuoteSymbol: msg.QuoteSymbol,
		Price:       msg.Price,
		Timestamp:   msg.Timestamp,
		BlockHeight: msg.BlockHeight,
		Sender:      msg.Sender,
	}
	b, _ := proto.Marshal(m)
	return b
}

func (d *RequestServer) broadcastVoteRequestPriceFeed(ctx sdktypes.Context, task *types.RequestPriceFeedIn, priceFeed *types.PriceFeedParam, rawPrices map[string]math.LegacyDec) error {
	d.Logger().Info("broadcastVoteRequestPriceFeed",
=======
	return s.aggregatePrices(pkgContext, request.Task.TaskIndex, request.Task.RequestId, priceFeedTxs)
}

func (s Server) broadcastVoteRequestPriceFeed(ctx sdktypes.Context, task *types.RequestPriceFeedIn, priceFeed *types.PriceFeedParam, rawPrices map[string]math.LegacyDec) error {
	s.Logger().Info("broadcastVoteRequestPriceFeed",
>>>>>>> main:dvs/price/server/msg_server_request_price_feed.go
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
		return fmt.Errorf("failed to broadcast VoteRequestPriceFeed for data error: %w", err)
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
<<<<<<< HEAD:dvs/price/server/request_server_price_feed.go
	bls := utils.SignWithBLS(d.blsKeyPair, d.getMsgBytes(&msg))
	msg.BlsSignature = bls

	if err := d.Server.SignAndBroadcastTx(ctx, &msg); err != nil {
		d.Logger().Error("broadcastVoteRequestPriceFeed SignAndBroadcastTx error: " + err.Error())
=======
	if err := s.SignAndBroadcastTx(ctx, &msg); err != nil {
		s.Logger().Error("broadcastVoteRequestPriceFeed SignAndBroadcastTx error: " + err.Error())
>>>>>>> main:dvs/price/server/msg_server_request_price_feed.go
		return fmt.Errorf("failed to broadcast VoteRequestPriceFeed for data error: %w", err)
	}

	return nil
}

<<<<<<< HEAD:dvs/price/server/request_server_price_feed.go
func (d *RequestServer) startCollectEnoughPriceFeedTxs(ctx sdktypes.Context, requestId []byte) ([]*pricetypes.MsgVoteRequestPriceFeed, error) {
	var (
		key       = string(requestId)
		eventData []*pricetypes.MsgVoteRequestPriceFeed
		blockData []*pricetypes.MsgVoteRequestPriceFeed
		ok        bool
	)

	// check if enough events
	events := d.PriceListener.QueryEvents(key)
	eventData, ok = d.verifyOperatorEvents(ctx, events)
	d.logger.Info("verifyOperatorEvents", "eventData", fmt.Sprintf("%+v", eventData), "ok", ok,
		"events length", len(events), "key", key)
	if ok {
		d.PriceListener.ClearEventsByKey(key)
		return eventData, nil
	}

	// get block height
	block, err := d.Server.GetLatestBlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block: %w", err)
	}
	blockHeight := block.Height

	// check if enough block height
	blocks := d.PriceListener.QueryBlocks(key)
	blockData, ok = d.checkAndChooseEnoughBlocks(ctx, blockHeight, blocks)
	d.logger.Info("checkAndChooseEnoughBlocks", "blockData", fmt.Sprintf("%+v", blockData), "ok", ok,
		"blockHeight", blockHeight, "blocks length", len(blocks), "key", key)
	if ok {
		d.PriceListener.ClearBlocksByKey(key)
		return blockData, nil
=======
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
>>>>>>> main:dvs/price/server/msg_server_request_price_feed.go
	}

	// subscribe events and blocks

	var (
		eventWaitChan      = make(chan struct{})
		blockWaitChan      = make(chan struct{})
		eventCh            = d.PriceListener.SubscribeEvents(1000)
		blockCh            = d.PriceListener.SubscribeBlocks(1000)
		operatorMaps       = make(map[string]*avsitypes.Operator)
		collectCtx, cancel = context.WithTimeout(ctx, 10*time.Minute)
	)
	for _, v := range ctx.Operators() {
		operatorMaps[string(v.Id)] = v
	}

	defer func() {
		cancel()

<<<<<<< HEAD:dvs/price/server/request_server_price_feed.go
		d.PriceListener.UnsubscribeEvents(eventCh)
		d.PriceListener.UnsubscribeBlocks(blockCh)

		close(eventWaitChan)
		close(blockWaitChan)

		d.PriceListener.ClearEventsByKey(key)
		d.PriceListener.ClearBlocksByKey(key)
	}()

	go d.collectEvents(collectCtx, operatorMaps, eventCh, &eventData, eventWaitChan) // listen events
	go d.collectBlocks(collectCtx, blockHeight, blockCh, &blockData, blockWaitChan)  // listen blocks

	select {
	case <-blockWaitChan:
		return blockData, nil
	case <-eventWaitChan:
		return eventData, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("failed to start collect enough price feed block timeout")
	}
}

func (d *RequestServer) verifyOperatorEvents(ctx sdktypes.Context, datas []tx_listener.EventData[*pricetypes.MsgVoteRequestPriceFeed]) ([]*pricetypes.MsgVoteRequestPriceFeed, bool) {
	var operatorMaps = make(map[string]*avsitypes.Operator)
	var operatorVerifyMaps = make(map[string]bool)
	for _, v := range ctx.Operators() {
		operatorMaps[string(v.Id)] = v
	}

	var out []*pricetypes.MsgVoteRequestPriceFeed
	for _, v := range datas {
		if operator, ok := operatorMaps[v.Data.OperatorId]; ok {
			// verify operator's key
			err := utils.VerifyBLSSignature(operator.Pubkeys, d.getMsgBytes(v.Data), v.Data.BlsSignature)
			if err != nil {
				continue
			}
			operatorVerifyMaps[string(operator.Id)] = true
			out = append(out, v.Data)
		}
	}

	return out, len(operatorVerifyMaps) == len(operatorMaps)
}

func (d *RequestServer) checkAndChooseEnoughBlocks(ctx context.Context, currentBlockHeight int64, blockDatas []tx_listener.BlockData[*pricetypes.MsgVoteRequestPriceFeed]) ([]*pricetypes.MsgVoteRequestPriceFeed, bool) {
	if len(blockDatas) == 0 {
		return nil, false
	}
	var out []*pricetypes.MsgVoteRequestPriceFeed
	sort.Slice(blockDatas, func(i, j int) bool {
		return blockDatas[i].Data.BlockHeight > blockDatas[j].Data.BlockHeight
	})
	var maxHeightBlocks = blockDatas[len(blockDatas)-1].Height

	for _, data := range blockDatas {
		if int64(data.Data.BlockHeight) > currentBlockHeight &&
			int64(data.Data.BlockHeight) <= currentBlockHeight+d.waitBlockCount {
			out = append(out, data.Data)
		}
	}

	return out, maxHeightBlocks >= currentBlockHeight+d.waitBlockCount
}

func (d *RequestServer) collectEvents(ctx context.Context, operatorMaps map[string]*avsitypes.Operator, ch *tx_listener.EventChannel[*pricetypes.MsgVoteRequestPriceFeed], eventData *[]*pricetypes.MsgVoteRequestPriceFeed, waitChan chan struct{}) {
	eventDataByOperatorId := make(map[string]*pricetypes.MsgVoteRequestPriceFeed)
	for _, data := range *eventData {
		eventDataByOperatorId[data.OperatorId] = data
	}

	mu := sync.Mutex{}
	select {
	case <-ctx.Done():
		return
	case <-ch.Done:
		return
	case data, ok := <-ch.Data:
		if !ok {
			return
		}
		// verify operator sign
		operator, ok := operatorMaps[data.Data.OperatorId]
		if ok && utils.VerifyBLSSignature(operator.Pubkeys, d.getMsgBytes(data.Data), data.Data.BlsSignature) == nil {
			mu.Lock()
			eventDataByOperatorId[data.Data.OperatorId] = data.Data
			mu.Unlock()
		}
		if len(eventDataByOperatorId) >= len(operatorMaps) {
			mu.Lock()
			// distinct by operator
			eventData = &[]*pricetypes.MsgVoteRequestPriceFeed{}
			for _, data := range eventDataByOperatorId {
				*eventData = append(*eventData, data)
			}
			mu.Unlock()
			waitChan <- struct{}{}
			return
		}
	}

}

func (d *RequestServer) collectBlocks(ctx context.Context, blockHeight int64, ch *tx_listener.BlockChannel[*pricetypes.MsgVoteRequestPriceFeed], blockData *[]*pricetypes.MsgVoteRequestPriceFeed, waitChan chan struct{}) {
	var mu = sync.Mutex{}
	select {
	case <-ctx.Done():
		return
	case <-ch.Done:
		return
	case data, ok := <-ch.Data:
		if !ok {
			return
		}
		if data.Height > blockHeight {
			mu.Lock()
			*blockData = append(*blockData, data.Data)
			mu.Unlock()
		}
		if data.Height >= blockHeight+d.waitBlockCount {
			waitChan <- struct{}{}
			return
		}
	}
}

func (d *RequestServer) aggregatePrices(ctx sdktypes.Context, taskIndex uint32, requestID []byte, priceFeedTxs []*pricetypes.MsgVoteRequestPriceFeed) (*types.RequestPriceFeedOut, error) {
=======
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
>>>>>>> main:dvs/price/server/msg_server_request_price_feed.go
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
