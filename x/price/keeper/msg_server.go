package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"intellix/x/price/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

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
		msg.Timestamp == 0 ||
		msg.BlockHeight == 0 {
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

func (k msgServer) VoteRequestPriceFeed(ctx context.Context, msg *types.MsgVoteRequestPriceFeed) (*types.MsgVoteRequestPriceFeedResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if err := k.validateMsgVoteRequestPriceFeed(msg); err != nil {
		return nil, err
	}

	if !k.IsOperatorHasPermission(sdkCtx, msg.OperatorId) {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "operator does not have permission to vote")
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

func (k msgServer) validateMsgVoteFinalizedRequestPrice(msg *types.MsgVoteFinalizedRequestPrice) error {
	if msg.TaskRaw == nil || msg.ValidatedData == nil || msg.PriceFeedResponse == nil {
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

	if err := store.Set(types.FinalizedRequestPrice(msg.TaskRaw.TaskIndex), data); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrIO, err.Error())
	}
	return nil
}

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
			sdk.NewAttribute(types.AttributeKeyTaskIndex, fmt.Sprintf("%d", msg.TaskRaw.TaskIndex)),
		),
	)

	return &types.MsgVoteFinalizedRequestPriceResponse{}, nil
}
