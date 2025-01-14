package tx_listener

import (
	"context"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
)

type EventHandler[K comparable, T any] func(ctx context.Context, event abci.Event) (K, T, error)
type BlockHandler[K comparable, T any] func(ctx context.Context, block *cmttypes.Block) (K, T, error)

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
