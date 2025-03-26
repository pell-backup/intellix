package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"intellix/x/price/types"
)

// QueryVoteRequestPriceFeed queries vote request price feed
func (k Keeper) QueryVoteRequestPriceFeed(goCtx context.Context, req *types.QueryVoteRequestPriceFeedReq) (*types.QueryVoteRequestPriceFeedResp, error) {
	if req == nil || req.RequestId == nil || req.OperatorId == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// query price feed vote
	store := k.storeService.OpenKVStore(ctx)
	data, err := store.Get(types.PriceFeedVoteKey(req.RequestId, req.OperatorId))
	if err != nil {
		return nil, status.Error(codes.NotFound, "price feed vote not found")
	}

	var priceFeed types.MsgVoteRequestPriceFeed
	err = k.cdc.Unmarshal(data, &priceFeed)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to unmarshal price feed vote")
	}

	var price []*types.VoteRequestPriceFeedPrice
	for _, p := range priceFeed.Price {
		price = append(price, &types.VoteRequestPriceFeedPrice{
			Source: p.Source,
			Price:  p.Price,
		})
	}

	return &types.QueryVoteRequestPriceFeedResp{
		TaskIndex:    priceFeed.TaskIndex,
		OperatorId:   priceFeed.OperatorId,
		RequestId:    priceFeed.RequestId,
		BaseSymbol:   priceFeed.BaseSymbol,
		QuoteSymbol:  priceFeed.QuoteSymbol,
		Price:        price,
		Timestamp:    priceFeed.Timestamp,
		BlockHeight:  priceFeed.BlockHeight,
		Sender:       priceFeed.Sender,
		BlsSignature: priceFeed.BlsSignature,
	}, nil
}
