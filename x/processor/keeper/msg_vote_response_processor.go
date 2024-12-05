package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"intellix/x/processor/types"
)

func (k msgServer) VoteResponseProcessor(goCtx context.Context, req *types.MsgVoteResponseProcessor) (*types.MsgVoteResponseProcessorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(goCtx)
	if err := k.validateMsgVoteResponseProcessor(req); err != nil {
		return nil, err
	}

	if err := k.saveMsgVoteResponseProcessor(sdkCtx, req); err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVoteResponseProcessor,
			sdk.NewAttribute(types.AttributeKeyTaskIndex, fmt.Sprintf("%d", req.TaskIndex)),
			sdk.NewAttribute(types.AttributeKeyRequestId, string(req.RequestId)),
		),
	)

	return &types.MsgVoteResponseProcessorResponse{}, nil
}

func (k msgServer) validateMsgVoteResponseProcessor(msg *types.MsgVoteResponseProcessor) error {
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
	if msg.ScriptId == 0 {
		return fmt.Errorf("script_id cannot be zero")
	}

	return nil
}

func (k msgServer) saveMsgVoteResponseProcessor(ctx sdk.Context, msg *types.MsgVoteResponseProcessor) error {
	store := k.storeService.OpenKVStore(ctx)
	data, err := k.cdc.Marshal(msg)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	if err := store.Set(types.MsgVoteResponseProcessorKey(msg.TaskIndex, msg.RequestId), data); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrIO, err.Error())
	}
	return nil
}
