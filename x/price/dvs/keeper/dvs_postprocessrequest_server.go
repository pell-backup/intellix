package keeper

import (
	"context"
	"intellix/x/price/dvs/types"
)

type DvsPostProcessRequestServer struct {
	Keeper
}

func NewDvsPostProcessRequestServer(keeper Keeper) types.DvsPostProcessRequestServer {
	return &DvsPostProcessRequestServer{Keeper: keeper}
}

func (d DvsPostProcessRequestServer) PostRequestPriceFeed(ctx context.Context, msg *types.PostPriceFeedMsg) (*types.ResponsePostRequestPriceFeed, error) {
	// TODO: add biz logic
	return nil, nil
}
