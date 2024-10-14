package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"intellix/x/price/dvs/types"
)

type DvsServer struct {
	Keeper
}

// NewDvsServerImpl returns an implementation of the DvsServer interface
// for the provided Keeper.
func NewDvsServerImpl(keeper Keeper) types.DvsServer {
	return &DvsServer{Keeper: keeper}
}

func (d *DvsServer) ProcessRequestPriceFeed(ctx context.Context, request *types.RequestProcessRequestPriceFeed) (*types.ResponseProcessRequestPriceFeed, error) {
	// new SDK context
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx = sdkCtx.WithBlockHeight(request.Height)
	sdkCtx = sdkCtx.WithChainID(request.ChainId.String())

	// TODO: add biz logic
	return nil, nil
}

func (d *DvsServer) PostRequestPriceFeed(ctx context.Context, request *types.RequestPostRequestPriceFeed) (*types.ResponsePostRequestPriceFeed, error) {
	// TODO: add biz logic
	return nil, nil
}
