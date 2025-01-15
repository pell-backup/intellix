package server

import (
	"context"
	"fmt"
	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	"github.com/cosmos/gogoproto/proto"
	"intellix/dvs/price/types"
	"intellix/pkg/tx_listener"
	"intellix/pkg/utils"
	sdktypes "intellix/sdk/types"
	pricetypes "intellix/x/price/types"
	"sync"
	"time"

	"cosmossdk.io/math"
	"sort"
)

func (d *RequestServer) RequestPriceFeed(ctx context.Context, request *types.RequestPriceFeedIn) (*types.RequestPriceFeedOut, error) {
	pkgContext := sdktypes.UnwrapContext(ctx)
	d.logger.Info("ProcessRequestPriceFeed", "PriceFeedParam", fmt.Sprintf("%+v", request.PriceFeed))

	// fetch raw price from chain
	rawPrices, err := fetchRawPrices(pkgContext, d.Logger(), request.PriceFeed.BaseSymbol, request.PriceFeed.QuoteSymbol)
	if err != nil {
		d.logger.Error("ProcessRequestPriceFeed fetchRawPrices error: " + err.Error())
		return nil, fmt.Errorf("failed to fetch raw prices: %w", err)
	}

	// sign data and broadcast VoteRequestPriceFeed
	height, err := d.broadcastVoteRequestPriceFeed(pkgContext, request, request.GetPriceFeed(), rawPrices)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast VoteRequestPriceFeed: %w", err)
	}

	priceFeedTxs, err := d.startCollectEnoughPriceFeedTxs(pkgContext, height, request.Task.RequestId)
	// listen and collect [N-N+M] block
	//priceFeedTxs, err := d.collectVoteRequestPriceFeedTxs(pkgContext, request.Task.RequestId)
	if err != nil {
		return nil, fmt.Errorf("failed to collect VoteRequestPriceFeed transactions: %w", err)
	}

	// aggregate [N-N+M] block prices
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

func (d *RequestServer) broadcastVoteRequestPriceFeed(ctx sdktypes.Context, task *types.RequestPriceFeedIn, priceFeed *types.PriceFeedParam, rawPrices map[string]math.LegacyDec) (int64, error) {
	d.Logger().Info("broadcastVoteRequestPriceFeed",
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
	addr, err := d.Server.SenderAddress()
	if err != nil {
		return 0, err
	}

	msg := pricetypes.MsgVoteRequestPriceFeed{
		Sender:      addr.String(),
		TaskIndex:   task.Task.TaskIndex,
		OperatorId:  d.Server.GetOperatorAddress(ctx),
		RequestId:   task.Task.RequestId,
		BaseSymbol:  priceFeed.BaseSymbol,
		QuoteSymbol: priceFeed.QuoteSymbol,
		Price:       prices,
		Timestamp:   time.Now().Unix(),
		BlockHeight: uint64(ctx.Height()),
	}
	bls := utils.SignWithBLS(d.blsKeyPair, d.getMsgBytes(&msg))
	msg.BlsSignature = bls

	height, err := d.Server.SignAndBroadcastTx(ctx, &msg)
	if err != nil {
		d.Logger().Error("broadcastVoteRequestPriceFeed SignAndBroadcastTx error: " + err.Error())
		return 0, fmt.Errorf("failed to broadcast VoteRequestPriceFeed for data error: %w", err)
	}

	return height, nil
}

func (d *RequestServer) startCollectEnoughPriceFeedTxs(ctx sdktypes.Context, blockHeight int64, requestId []byte) ([]*pricetypes.MsgVoteRequestPriceFeed, error) {
	key := string(requestId)

	collectCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	var (
		eventData []*pricetypes.MsgVoteRequestPriceFeed
		blockData []*pricetypes.MsgVoteRequestPriceFeed
		ok        bool
	)

	// check if enough events
	events := d.PriceListener.QueryEvents(key)
	eventData, ok = d.verifyOperatorEvents(ctx, events)
	if ok {
		d.PriceListener.ClearEventsByKey(key)
		return eventData, nil
	}

	// check if enough block height
	blocks := d.PriceListener.QueryBlocks(key)
	blockData, ok = d.checkAndChooseEnoughBlocks(ctx, blockHeight, blocks)
	if ok {
		d.PriceListener.ClearBlocksByKey(key)
		return blockData, nil
	}

	var (
		eventWaitChan = make(chan struct{})
		blockWaitChan = make(chan struct{})
		operatorMaps  = make(map[string]*avsitypes.Operator)
	)
	defer func() {
		close(eventWaitChan)
		close(blockWaitChan)
	}()

	for _, v := range ctx.Operators() {
		operatorMaps[string(v.Id)] = v
	}

	var (
		eventCh = d.PriceListener.SubscribeEvents(1000)
		blockCh = d.PriceListener.SubscribeBlocks(1000)
	)
	defer func() {
		d.PriceListener.UnsubscribeEvents(eventCh)
		d.PriceListener.UnsubscribeBlocks(blockCh)
	}()

	go d.collectEvents(collectCtx, operatorMaps, eventCh, eventData, eventWaitChan) // listen events
	go d.collectBlocks(collectCtx, blockHeight, blockCh, blockData, blockWaitChan)  // listen blocks

	select {
	case <-blockWaitChan:
		return blockData, nil
	case <-eventWaitChan:
		return blockData, nil
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

func (d *RequestServer) collectEvents(ctx context.Context, operatorMaps map[string]*avsitypes.Operator, ch *tx_listener.EventChannel[*pricetypes.MsgVoteRequestPriceFeed], eventData []*pricetypes.MsgVoteRequestPriceFeed, waitChan chan struct{}) {
	eventDataByOperatorId := make(map[string]*pricetypes.MsgVoteRequestPriceFeed)
	for _, data := range eventData {
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
			eventData = []*pricetypes.MsgVoteRequestPriceFeed{}
			for _, data := range eventDataByOperatorId {
				eventData = append(eventData, data)
			}
			mu.Unlock()
			waitChan <- struct{}{}
			return
		}
	}

}

func (d *RequestServer) collectBlocks(ctx context.Context, blockHeight int64, ch *tx_listener.BlockChannel[*pricetypes.MsgVoteRequestPriceFeed], blockData []*pricetypes.MsgVoteRequestPriceFeed, waitChan chan struct{}) {
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
			blockData = append(blockData, data.Data)
			mu.Unlock()
		}
		if data.Height >= blockHeight+d.waitBlockCount {
			waitChan <- struct{}{}
			return
		}
	}
}

func (d *RequestServer) aggregatePrices(ctx sdktypes.Context, taskIndex uint32, requestID []byte, priceFeedTxs []*pricetypes.MsgVoteRequestPriceFeed) (*types.RequestPriceFeedOut, error) {
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
