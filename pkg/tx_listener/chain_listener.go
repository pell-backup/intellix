package tx_listener

import (
	"context"
	"github.com/0xPellNetwork/pelldvs/libs/log"
	"github.com/cosmos/cosmos-sdk/client"
	"sync"
	"time"

	cmttypes "github.com/cometbft/cometbft/types"

	tmclient "github.com/cometbft/cometbft/rpc/client/http"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type ChainListenerIFace[E any, B any] interface {
	Start() error
	Stop()

	QueryEvents() []EventData[E]
	QueryBlocks(height int64) []BlockData[B]

	SubscribeEvents(maxQueueSize int) *EventChannel[E]
	SubscribeBlocks(maxQueueSize int) *BlockChannel[B]

	ClearEvents()
	ClearBlocks()

	ClearEventsByFunc(func(EventData[E]) bool)
	ClearBlocksByFunc(func(BlockData[B]) bool)
}

type ChainListener[E any, B any] struct {
	sync.RWMutex
	logger log.Logger

	ctx    context.Context
	cancel context.CancelFunc

	clientCtx      client.Context
	wsEndpoint     string // websocket endpoint
	subscribeQuery string // event query string

	eventHandler EventHandler[E] // new event handler
	blockHandler BlockHandler[B] // new block handler

	events     []EventData[E]
	blocks     []BlockData[B]
	maxHistory int // max event or block size

	// listeners
	eventChannels []*EventChannel[E]
	blockChannels []*BlockChannel[B]
}

func NewChainListener[E any, B any](
	logger log.Logger,
	clientCtx client.Context,
	wsEndpoint string,
	subscribeQuery string,
	maxHistory int,
	eventHandler EventHandler[E],
	blockHandler BlockHandler[B],
) *ChainListener[E, B] {
	ctx, cancel := context.WithCancel(context.Background())
	if maxHistory <= 0 {
		maxHistory = 1000 // default
	}
	if eventHandler == nil || blockHandler == nil {
		panic("handlers cannot be nil")
	}
	return &ChainListener[E, B]{
		logger:         logger,
		ctx:            ctx,
		cancel:         cancel,
		clientCtx:      clientCtx,
		wsEndpoint:     wsEndpoint,
		subscribeQuery: subscribeQuery,
		maxHistory:     maxHistory,
		eventHandler:   eventHandler,
		blockHandler:   blockHandler,
		events:         make([]EventData[E], 0),
		blocks:         make([]BlockData[B], 0),
		eventChannels:  make([]*EventChannel[E], 0),
		blockChannels:  make([]*BlockChannel[B], 0),
	}
}

func (l *ChainListener[E, B]) Start() error {
	l.Lock()
	defer l.Unlock()

	go l.startWebsocketListener(l.ctx)
	go l.startBlockScanner(l.ctx)
	return nil
}

func (l *ChainListener[E, B]) Stop() {
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

func (l *ChainListener[E, B]) startWebsocketListener(ctx context.Context) {
	cli, err := tmclient.New(l.wsEndpoint, "/websocket")
	if err != nil {
		return
	}

	if err := cli.Start(); err != nil {
		return
	}
	defer cli.Stop()

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

func (l *ChainListener[E, B]) startBlockScanner(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			node, err := l.clientCtx.GetNode()

			block, err := node.Block(ctx, nil)
			if err != nil {
				continue
			}
			l.handleNewBlock(ctx, block.Block)
		}
	}
}

func (l *ChainListener[E, B]) QueryEvents() []EventData[E] {
	l.RLock()
	defer l.RUnlock()

	result := make([]EventData[E], len(l.events))
	copy(result, l.events)
	return result
}

func (l *ChainListener[E, B]) QueryBlocks(height int64) []BlockData[B] {
	l.RLock()
	defer l.RUnlock()

	var results []BlockData[B]
	for _, block := range l.blocks {
		if block.Height >= height {
			results = append(results, block)
		}
	}
	return results
}

func (l *ChainListener[E, B]) SubscribeEvents(maxQueueSize int) *EventChannel[E] {
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

func (l *ChainListener[E, B]) SubscribeBlocks(maxQueueSize int) *BlockChannel[B] {
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

func (l *ChainListener[E, B]) ClearEvents() {
	l.Lock()
	l.events = make([]EventData[E], 0)
	l.Unlock()
}

func (l *ChainListener[E, B]) ClearBlocks() {
	l.Lock()
	l.blocks = make([]BlockData[B], 0)
	l.Unlock()
}

func (l *ChainListener[E, B]) ClearEventsByFunc(f func(EventData[E]) bool) {
	l.Lock()
	defer l.Unlock()

	newEvents := make([]EventData[E], 0, len(l.events))
	for _, event := range l.events {
		if !f(event) {
			newEvents = append(newEvents, event)
		}
	}
	l.events = newEvents
}

func (l *ChainListener[E, B]) ClearBlocksByFunc(f func(BlockData[B]) bool) {
	l.Lock()
	defer l.Unlock()

	newBlocks := make([]BlockData[B], 0, len(l.blocks))
	for _, block := range l.blocks {
		if !f(block) {
			newBlocks = append(newBlocks, block)
		}
	}
	l.blocks = newBlocks
}

func (l *ChainListener[E, B]) handleNewEvent(ctx context.Context, tmEvent tmctypes.ResultEvent) {
	if txData, ok := tmEvent.Data.(cmttypes.EventDataTx); ok {
		// handle tx
		for _, event := range txData.Result.Events {
			if txData.Result.Code != 0 {
				l.logger.Error("tx failed", "code", txData.Result.Code, "event", event, "hash", txData.Tx)
				continue
			}

			// invoke event handler
			eventData, err := l.eventHandler(ctx, event)
			if err != nil {
				l.logger.Error("event handler failed", "err", err, "event", event, "hash", txData.Tx)
				continue
			}

			ed := EventData[E]{
				Height:    txData.Height,
				TxHash:    txData.Tx,
				Event:     event,
				Data:      eventData,
				Timestamp: time.Now(),
			}
			// save to memory
			l.saveAndBroadcastEvent(ed)
		}
	}
}

func (l *ChainListener[E, B]) saveAndBroadcastEvent(eventData EventData[E]) {
	l.Lock()
	l.events = append(l.events, eventData)
	// remove old events
	if len(l.events) > l.maxHistory {
		l.events = l.events[1:]
	}
	l.Unlock()

	// broadcast to all subscribers
	for _, channel := range l.eventChannels {
		select {
		case channel.Data <- eventData:
		default:
			// no block if channel is full
		}
	}
}

func (l *ChainListener[E, B]) handleNewBlock(ctx context.Context, block *cmttypes.Block) {
	data, err := l.blockHandler(ctx, block)
	if err != nil {
		return
	}

	blockData := BlockData[B]{
		Height:    block.Height,
		Block:     block,
		Data:      data,
		Timestamp: time.Now(),
	}

	l.Lock()
	l.blocks = append(l.blocks, blockData)
	// remove old blocks
	if len(l.blocks) > l.maxHistory {
		l.blocks = l.blocks[1:]
	}
	l.Unlock()

	for _, channel := range l.blockChannels {
		select {
		case channel.Data <- blockData:
		default:
			// no block if channel is full
		}
	}
}
