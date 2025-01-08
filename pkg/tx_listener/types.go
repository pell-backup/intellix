package tx_listener

import (
	"context"
	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"time"
)

type EventHandler[T any] func(ctx context.Context, event abci.Event) (T, error)
type BlockHandler[T any] func(ctx context.Context, block *cmttypes.Block) (T, error)

type EventData[T any] struct {
	Height    int64
	TxHash    []byte
	Event     abci.Event
	Data      T // data after processing
	Timestamp time.Time
}

type BlockData[T any] struct {
	Height    int64
	Block     *cmttypes.Block
	Data      T // data after processing
	Timestamp time.Time
}

type EventChannel[T any] struct {
	Data chan EventData[T]
	Done chan struct{}
}

type BlockChannel[T any] struct {
	Data chan BlockData[T]
	Done chan struct{}
}
