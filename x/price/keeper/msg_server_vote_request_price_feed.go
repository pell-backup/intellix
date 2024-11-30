package keeper

import (
	"context"
	"fmt"
	"intellix/x/price/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) VoteRequestPriceFeed(ctx context.Context, msg *types.MsgVoteRequestPriceFeed) (*types.MsgVoteRequestPriceFeedResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.validateMsgVoteRequestPriceFeed(msg); err != nil {
		return nil, err
	}

	if err := k.savePriceFeedVote(sdkCtx, msg); err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVoteRequestPriceFeed,
			sdk.NewAttribute(types.AttributeKeyTaskIndex, fmt.Sprintf("%d", msg.TaskIndex)),
			sdk.NewAttribute(types.AttributeKeyOperatorId, msg.OperatorId),
			sdk.NewAttribute(types.AttributeKeyRequestId, string(msg.RequestId)),
		),
	)

	return &types.MsgVoteRequestPriceFeedResponse{}, nil
}

func (k msgServer) validateMsgVoteRequestPriceFeed(msg *types.MsgVoteRequestPriceFeed) error {
	if msg == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "msg cannot be nil")
	}
	if msg.TaskIndex == 0 ||
		msg.OperatorId == "" ||
		msg.RequestId == nil ||
		msg.BaseSymbol == "" ||
		msg.QuoteSymbol == "" ||
		msg.Price == nil ||
		msg.Timestamp == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid msg")
	}

	return nil
}

func (k msgServer) savePriceFeedVote(ctx sdk.Context, msg *types.MsgVoteRequestPriceFeed) error {
	store := k.storeService.OpenKVStore(ctx)
	data, err := k.cdc.Marshal(msg)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	if err := store.Set(types.PriceFeedVoteKey(msg.TaskIndex, msg.OperatorId), data); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrIO, err.Error())
	}
	return nil
}
