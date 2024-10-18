package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	"intellix/x/price/dvs/types"
)

type DvsPostProcessRequestServer struct {
	Keeper
}

func NewDvsPostProcessRequestServer(keeper Keeper) types.DvsPostProcessRequestServer {
	return &DvsPostProcessRequestServer{Keeper: keeper}
}

func (d DvsPostProcessRequestServer) PostProcessRequestPriceFeed(ctx context.Context, msg *types.RequestPostRequestPriceFeedValidatedData) (*types.ResponsePostRequestPriceFeed, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	aggrMsg, err := dvsservermanager.DecodeMsg(msg.Data)
	if err != nil {
		return nil, err
	}
	aggrData, ok := aggrMsg.(*types.AggregatedRequestPrice)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &types.AggregatedRequestPrice{}, aggrMsg)
	}

	// just send VoteFinalizedRequestPrice Tx
	err = d.sendVoteFinalizedRequestPriceTx(sdkCtx, msg, aggrData)
	if err != nil {
		return nil, err
	}

	return &types.ResponsePostRequestPriceFeed{}, nil
}

func (d DvsPostProcessRequestServer) sendVoteFinalizedRequestPriceTx(ctx sdk.Context, raw *types.RequestPostRequestPriceFeedValidatedData, priceData *types.AggregatedRequestPrice) error {
	if err := d.Keeper.SignAndBroadcastTx(ctx, priceData); err != nil {
		return err
	}
	return nil
}
