package keeper

import (
	"context"
	"fmt"
	"intellix/x/price/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) VoteFinalizedRequestPrice(ctx context.Context, msg *types.MsgVoteFinalizedRequestPrice) (*types.MsgVoteFinalizedRequestPriceResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.validateMsgVoteFinalizedRequestPrice(msg); err != nil {
		return nil, err
	}
	if err := k.saveFinalizedRequestPrice(sdkCtx, msg); err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVoteFinalizedRequestPrice,
			sdk.NewAttribute(types.AttributeKeyTaskIndex, fmt.Sprintf("%d", msg.TaskIndex)),
		),
	)

	return &types.MsgVoteFinalizedRequestPriceResponse{}, nil
}

func (k msgServer) validateMsgVoteFinalizedRequestPrice(msg *types.MsgVoteFinalizedRequestPrice) error {
	if msg.RequestId == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid msg")
	}

	return nil
}

func (k msgServer) saveFinalizedRequestPrice(ctx sdk.Context, msg *types.MsgVoteFinalizedRequestPrice) error {
	store := k.storeService.OpenKVStore(ctx)
	data, err := k.cdc.Marshal(msg)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	if err := store.Set(types.FinalizedRequestPrice(msg.TaskIndex), data); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrIO, err.Error())
	}
	return nil
}
