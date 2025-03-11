package tx_listener

import (
	"context"
	"sync"
	"time"

	"github.com/0xPellNetwork/pelldvs-libs/log"
	"github.com/cosmos/cosmos-sdk/client"

	cmttypes "github.com/cometbft/cometbft/types"

	tmclient "github.com/cometbft/cometbft/rpc/client/http"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type ChainListenerIFace[K comparable, E any, B any] interface {
	Start()
	Stop()

	QueryEvents(key K) []EventData[E]
	QueryBlocks(key K) []BlockData[B]

	SubscribeEvents(maxQueueSize int) *EventChannel[E]
	SubscribeBlocks(maxQueueSize int) *BlockChannel[B]

	UnsubscribeEvents(ch *EventChannel[E])
	UnsubscribeBlocks(ch *BlockChannel[B])

	ClearEvents()
	ClearBlocks()

	ClearEventsByKey(key K)
	ClearBlocksByKey(key K)
}

type ChainListener[K comparable, E any, B any] struct {
	sync.RWMutex
	logger log.Logger

	ctx    context.Context
	cancel context.CancelFunc

	clientCtx      client.Context
	wsEndpoint     string // websocket endpoint
	subscribeQuery string // event query string

	eventHandler EventHandler[K, E] // new event handler
	blockHandler BlockHandler[K, B] // new block handler

	// maybe save to db
	events     map[K][]EventData[E]
	blocks     map[K][]BlockData[B]
	maxHistory int // max event or block size

	cleanupInterval time.Duration // clean data interval
	retentionPeriod time.Duration // data retent time

	// listeners
	eventChannels []*EventChannel[E]
	blockChannels []*BlockChannel[B]
}

func NewChainListener[K comparable, E any, B any](
	logger log.Logger,
	clientCtx client.Context,
	wsEndpoint string,
	subscribeQuery string,
	maxHistory int,
	eventHandler EventHandler[K, E],
	blockHandler BlockHandler[K, B],
) *ChainListener[K, E, B] {
	ctx, cancel := context.WithCancel(context.Background())
	if maxHistory <= 0 {
		maxHistory = 1000 // default
	}
	if eventHandler == nil || blockHandler == nil {
		panic("handlers cannot be nil")
	}
	return &ChainListener[K, E, B]{
		logger:          logger,
		ctx:             ctx,
		cancel:          cancel,
		clientCtx:       clientCtx,
		wsEndpoint:      wsEndpoint,
		subscribeQuery:  subscribeQuery,
		maxHistory:      maxHistory,
		eventHandler:    eventHandler,
		blockHandler:    blockHandler,
		events:          make(map[K][]EventData[E]),
		blocks:          make(map[K][]BlockData[B]),
		eventChannels:   make([]*EventChannel[E], 0),
		blockChannels:   make([]*BlockChannel[B], 0),
		cleanupInterval: 30 * time.Minute,
		retentionPeriod: 60 * time.Minute,
	}
}

func (l *ChainListener[K, E, B]) Start() {
	l.Lock()
	defer l.Unlock()

	go l.startWebsocketListener(l.ctx)
	go l.startBlockScanner(l.ctx)
	go l.startCleanupTask(l.ctx)
}

func (l *ChainListener[K, E, B]) Stop() {
	l.Lock()
	defer l.Unlock()

	l.cancel()
	for _, ch := range l.eventChannels {
		close(ch.Done)
		close(ch.Data)
	}
	l.eventChannels = nil

	for _, ch := range l.blockChannels {
		close(ch.Done)
		close(ch.Data)
	}
	l.blockChannels = nil

	l.events = nil
	l.blocks = nil
}

func (l *ChainListener[K, E, B]) startWebsocketListener(ctx context.Context) {
	cli, err := tmclient.New(l.wsEndpoint, "/websocket")
	if err != nil {
		return
	}

	if err := cli.Start(); err != nil {
		return
	}
	defer func() {
		_ = cli.Stop()
	}()

	eventCh, err := cli.Subscribe(ctx, "chain-listener", l.subscribeQuery)
	if err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case tmEvent := <-eventCh:
			l.handleNewEvent(ctx, tmEvent)
		}
	}
}

func (l *ChainListener[K, E, B]) startBlockScanner(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			node, err := l.clientCtx.GetNode()
			if err != nil {
				continue
			}

			block, err := node.Block(ctx, nil)
			if err != nil {
				continue
			}
			l.handleNewBlock(ctx, block.Block)
		}
	}
}

func (l *ChainListener[K, E, B]) startCleanupTask(ctx context.Context) {
	ticker := time.NewTicker(l.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.cleanup()
		}
	}
}

func (l *ChainListener[K, E, B]) cleanup() {
	l.Lock()
	defer l.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.retentionPeriod)

	for key, events := range l.events {
		var newEvents []EventData[E]
		for _, event := range events {
			if event.Timestamp.After(cutoff) {
				newEvents = append(newEvents, event)
			}
		}
		if len(newEvents) > 0 {
			l.events[key] = newEvents
		} else {
			delete(l.events, key)
		}
	}

	for key, blocks := range l.blocks {
		var newBlocks []BlockData[B]
		for _, block := range blocks {
			if block.Timestamp.After(cutoff) {
				newBlocks = append(newBlocks, block)
			}
		}
		if len(newBlocks) > 0 {
			l.blocks[key] = newBlocks
		} else {
			delete(l.blocks, key)
		}
	}

	l.logger.Info("Cleanup completed",
		"events_count", len(l.events),
		"blocks_count", len(l.blocks),
		"cutoff_time", cutoff,
	)
}

func (l *ChainListener[K, E, B]) QueryEvents(key K) []EventData[E] {
	l.RLock()
	defer l.RUnlock()

	return l.events[key]
}

func (l *ChainListener[K, E, B]) QueryBlocks(key K) []BlockData[B] {
	l.RLock()
	defer l.RUnlock()

	return l.blocks[key]
}

func (l *ChainListener[K, E, B]) SubscribeEvents(maxQueueSize int) *EventChannel[E] {
	if maxQueueSize <= 0 {
		maxQueueSize = 100
	}

	channel := &EventChannel[E]{
		Data: make(chan EventData[E], maxQueueSize),
		Done: make(chan struct{}),
	}

	l.Lock()
	l.eventChannels = append(l.eventChannels, channel)
	l.Unlock()

	return channel
}

func (l *ChainListener[K, E, B]) SubscribeBlocks(maxQueueSize int) *BlockChannel[B] {
	if maxQueueSize <= 0 {
		maxQueueSize = 100
	}

	channel := &BlockChannel[B]{
		Data: make(chan BlockData[B], maxQueueSize),
		Done: make(chan struct{}),
	}

	l.Lock()
	l.blockChannels = append(l.blockChannels, channel)
	l.Unlock()

	return channel
}

func (l *ChainListener[K, E, B]) UnsubscribeEvents(ch *EventChannel[E]) {
	l.Lock()
	defer l.Unlock()

	close(ch.Done)
	close(ch.Data)

	for i, channel := range l.eventChannels {
		if channel == ch {
			l.eventChannels = append(l.eventChannels[:i], l.eventChannels[i+1:]...)
			break
		}
	}
}

func (l *ChainListener[K, E, B]) UnsubscribeBlocks(ch *BlockChannel[B]) {
	l.Lock()
	defer l.Unlock()

	close(ch.Done)
	close(ch.Data)

	for i, channel := range l.blockChannels {
		if channel == ch {
			l.blockChannels = append(l.blockChannels[:i], l.blockChannels[i+1:]...)
			break
		}
	}
}

func (l *ChainListener[K, E, B]) ClearEvents() {
	l.Lock()
	l.events = make(map[K][]EventData[E])
	l.Unlock()
}

func (l *ChainListener[K, E, B]) ClearBlocks() {
	l.Lock()
	l.blocks = make(map[K][]BlockData[B])
	l.Unlock()
}

func (l *ChainListener[K, E, B]) ClearEventsByKey(key K) {
	l.Lock()
	defer l.Unlock()

	delete(l.events, key)
}

func (l *ChainListener[K, E, B]) ClearBlocksByKey(key K) {
	l.Lock()
	defer l.Unlock()

	delete(l.blocks, key)
}

func (l *ChainListener[K, E, B]) handleNewEvent(ctx context.Context, tmEvent tmctypes.ResultEvent) {
	if txData, ok := tmEvent.Data.(cmttypes.EventDataTx); ok {
		for _, event := range txData.Result.Events {
			if txData.Result.Code != 0 {
				l.logger.Error("tx failed", "code", txData.Result.Code)
				continue
			}

			key, eventData, err := l.eventHandler(ctx, event)
			if err != nil {
				continue
			}

			var zero K
			if key == zero {
				continue
			}

			l.logger.Info("ChainListener.eventHandler", "key", key, "event", event, "eventData", eventData)

			ed := EventData[E]{
				Height:    txData.Height,
				TxHash:    txData.Tx,
				Event:     event,
				Data:      eventData,
				Timestamp: time.Now(),
			}

			l.saveAndBroadcastEvent(key, ed)
		}
	}
}

func (l *ChainListener[K, E, B]) saveAndBroadcastEvent(key K, eventData EventData[E]) {
	l.Lock()
	if len(l.events) >= l.maxHistory {
		var oldestKey K
		var oldestTime time.Time
		first := true
		for _, v := range l.events[key] {
			if first || v.Timestamp.Before(oldestTime) {
				oldestKey = key
				oldestTime = v.Timestamp
				first = false
			}
		}
		delete(l.events, oldestKey)
	}
	l.events[key] = append(l.events[key], eventData)
	l.Unlock()

	for _, channel := range l.eventChannels {
		select {
		case <-channel.Done:
			continue
		default:
			select {
			case channel.Data <- eventData:
			default:
				// not block
			}
		}
	}
}

func (l *ChainListener[K, E, B]) handleNewBlock(ctx context.Context, block *cmttypes.Block) {
	key, data, err := l.blockHandler(ctx, block)
	if err != nil {
		//l.logger.Error("blockHandler", "error", err)
		return
	}
	var zero K
	if key == zero {
		return
	}

	l.logger.Info("ChainListener.blockHandler", "key", key, "data", data)

	blockData := BlockData[B]{
		Height:    block.Height,
		Block:     block,
		Data:      data,
		Timestamp: time.Now(),
	}

	l.Lock()
	if len(l.blocks[key]) >= l.maxHistory {
		var oldestKey K
		var oldestTime time.Time
		first := true
		for _, v := range l.blocks[key] {
			if first || v.Timestamp.Before(oldestTime) {
				oldestKey = key
				oldestTime = v.Timestamp
				first = false
			}
		}
		delete(l.blocks, oldestKey)
	}
	l.blocks[key] = append(l.blocks[key], blockData)
	l.Unlock()

	for _, channel := range l.blockChannels {
		select {
		case <-channel.Done:
			continue
		default:
			select {
			case channel.Data <- blockData:
			default:
				// not block
			}
		}
	}
}
