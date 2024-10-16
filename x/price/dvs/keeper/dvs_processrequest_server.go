package keeper

import (
	"context"
	"intellix/x/price/dvs/types"
)

type DvsProcessRequestServer struct {
	Keeper
}

// NewDvsProcessRequestServer returns an implementation of the DvsProcessRequestServer interface
// for the provided Keeper.
func NewDvsProcessRequestServer(keeper Keeper) types.DvsProcessRequestServer {
	return &DvsProcessRequestServer{Keeper: keeper}
}

func (d *DvsProcessRequestServer) ProcessRequestPriceFeed(ctx context.Context, request *types.ProcessPriceFeedMsg) (*types.ProcessPriceFeedResp, error) {
	// TODO: add biz logic
	return nil, nil
}
