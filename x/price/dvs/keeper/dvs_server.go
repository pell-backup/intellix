package keeper

import (
	"context"
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

func (d *DvsServer) ProcessRequestPriceFeed(ctx context.Context, request *types.ProcessPriceFeedMsg) (*types.ProcessPriceFeedResp, error) {
	// TODO: add biz logic
	return nil, nil
}

func (d *DvsServer) PostRequestPriceFeed(ctx context.Context, request *types.PostPriceFeedMsg) (*types.ResponsePostRequestPriceFeed, error) {
	// TODO: add biz logic
	return nil, nil
}
