package keeper

import (
	"context"
	"fmt"
	"intellix/x/processor/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) VoteRequestProcessor(goCtx context.Context, req *types.MsgVoteRequestProcessor) (*types.MsgVoteRequestProcessorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(goCtx)

	if err := k.validateMsgVoteRequestProcessor(req); err != nil {
		return nil, err
	}

	if err := k.saveMsgVoteRequestProcessor(sdkCtx, req); err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVoteRequestProcessor,
			sdk.NewAttribute(types.AttributeKeyTaskIndex, fmt.Sprintf("%d", req.TaskIndex)),
			sdk.NewAttribute(types.AttributeKeyOperatorId, req.OperatorId),
			sdk.NewAttribute(types.AttributeKeyRequestId, string(req.RequestId)),
		),
	)

	return &types.MsgVoteRequestProcessorResponse{}, nil
}

func (k msgServer) validateMsgVoteRequestProcessor(msg *types.MsgVoteRequestProcessor) error {
	if msg.Sender == "" {
		return fmt.Errorf("sender cannot be empty")
	}
	if msg.TaskIndex == 0 {
		return fmt.Errorf("task_index cannot be zero")
	}
	if msg.RequestId == nil {
		return fmt.Errorf("request_id cannot be empty")
	}
	if msg.CallbackFunctionId == nil {
		return fmt.Errorf("callback_function_id cannot be empty")
	}
	if msg.CallbackAddress == "" {
		return fmt.Errorf("callback_address cannot be empty")
	}

	return nil
}

func (k msgServer) saveMsgVoteRequestProcessor(ctx sdk.Context, msg *types.MsgVoteRequestProcessor) error {
	store := k.storeService.OpenKVStore(ctx)
	data, err := k.cdc.Marshal(msg)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	if err := store.Set(types.MsgVoteRequestProcessorKey(msg.RequestId, msg.OperatorId), data); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrIO, err.Error())
	}
	return nil
}
