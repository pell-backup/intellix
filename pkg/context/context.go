package pkgcontext

import (
	"context"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/gogoproto/proto"
	"time"
)

type ContextKeyType string

const (
	ContextKey ContextKeyType = "pkg_context"
)

type Context struct {
	baseCtx context.Context

	chainID string
	header  cmtproto.Header
}

func NewContext(baseCtx context.Context, header *cmtproto.Header) Context {
	if header == nil {
		header = &cmtproto.Header{}
	}

	return Context{
		baseCtx: baseCtx,
		header:  *header,
	}
}

func (c Context) Deadline() (deadline time.Time, ok bool) {
	return c.baseCtx.Deadline()
}

func (c Context) Done() <-chan struct{} {
	return c.baseCtx.Done()
}

func (c Context) Err() error {
	return c.baseCtx.Err()
}

func (c Context) Value(key any) any {
	if key == ContextKey {
		return c
	}

	return c.baseCtx.Value(key)
}

func (c Context) WithValue(key, value interface{}) Context {
	c.baseCtx = context.WithValue(c.baseCtx, key, value)
	return c
}

func (c Context) BlockHeight() int64 { return c.header.Height }

func (c Context) WithBlockHeight(height int64) Context {
	newHeader := c.BlockHeader()
	newHeader.Height = height
	return c.WithBlockHeader(newHeader)
}

func (c Context) BlockTime() time.Time { return c.header.Time }

// WithBlockTime returns a Context with an updated CometBFT block header time in UTC with no monotonic component.
// Stripping the monotonic component is for time equality.
func (c Context) WithBlockTime(newTime time.Time) Context {
	newHeader := c.BlockHeader()
	// https://github.com/gogo/protobuf/issues/519
	newHeader.Time = newTime.Round(0).UTC()
	return c.WithBlockHeader(newHeader)
}

func (c Context) BlockHeader() cmtproto.Header {
	msg := proto.Clone(&c.header).(*cmtproto.Header)
	return *msg
}

func (c Context) ChainID() string {
	return c.chainID
}

func (c Context) WithContext(ctx context.Context) Context {
	c.baseCtx = ctx
	return c
}

// WithBlockHeader returns a Context with an updated CometBFT block header in UTC time.
func (c Context) WithBlockHeader(header cmtproto.Header) Context {
	// https://github.com/gogo/protobuf/issues/519
	header.Time = header.Time.UTC()
	c.header = header
	return c
}

// WithChainID returns a Context with an updated chain identifier.
func (c Context) WithChainID(chainID string) Context {
	c.chainID = chainID
	return c
}

func UnwrapContext(ctx context.Context) Context {
	if sdkCtx, ok := ctx.(Context); ok {
		return sdkCtx
	}
	return ctx.Value(ContextKey).(Context)
}
