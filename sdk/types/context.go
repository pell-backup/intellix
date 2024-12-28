package types

import (
	"context"
	"time"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
)

type ContextKeyType string

const (
	ContextKey ContextKeyType = "pkg_context"
)

type Context struct {
	baseCtx                   context.Context
	header                    cmtproto.Header
	chainID                   int64
	height                    int64
	groupNumbers              []uint32
	groupThresholdPercentages []uint32
	dvsPostResponseData       []byte
	//TODO  Operators
}

// Read-only accessors
func (c Context) Context() context.Context            { return c.baseCtx }
func (c Context) ChainID() int64                      { return c.chainID }
func (c Context) Height() int64                       { return c.height }
func (c Context) GroupNumbers() []uint32              { return c.groupNumbers }
func (c Context) GroupThresholdPercentages() []uint32 { return c.groupThresholdPercentages }
func (c Context) DvsPostResponseData() []byte         { return c.dvsPostResponseData }

func (c Context) Value(key any) any {
	if key == ContextKey {
		return c
	}
	return c.baseCtx.Value(key)
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

func NewContext(baseCtx context.Context, header *cmtproto.Header) Context {
	if header == nil {
		header = &cmtproto.Header{}
	}

	return Context{
		baseCtx: baseCtx,
		header:  *header,
	}
}

func (c Context) WithValue(key, value interface{}) Context {
	c.baseCtx = context.WithValue(c.baseCtx, key, value)
	return c
}

// WithContext returns a Context with an updated context.Context.
func (c Context) WithContext(ctx context.Context) Context {
	c.baseCtx = ctx
	return c
}

func (c Context) WithChainID(chainID int64) Context {
	c.chainID = chainID
	return c
}

func (c Context) WithHeight(height int64) Context {
	c.height = height
	return c
}

func (c Context) WithGroupNumbers(groupNumbers []uint32) Context {
	c.groupNumbers = groupNumbers
	return c
}

func (c Context) WithGroupThresholdPercentages(groupThresholdPercentages []uint32) Context {
	c.groupThresholdPercentages = groupThresholdPercentages
	return c
}

// WithBlockHeader returns a Context with an updated CometBFT block header in UTC time.
func (c Context) WithBlockHeader(header cmtproto.Header) Context {
	// https://github.com/gogo/protobuf/issues/519
	header.Time = header.Time.UTC()
	c.header = header
	return c
}

func (c Context) WithDvsPostResponseData(postProcessResponseData []byte) Context {
	c.dvsPostResponseData = postProcessResponseData
	return c
}

func UnwrapContext(ctx context.Context) Context {
	if sdkCtx, ok := ctx.(Context); ok {
		return sdkCtx
	}
	return ctx.Value(ContextKey).(Context)
}
