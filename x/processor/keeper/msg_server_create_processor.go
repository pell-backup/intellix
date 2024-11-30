package keeper

import (
	"context"
	"intellix/x/processor/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreateProcessor(goCtx context.Context, msg *types.MsgCreateProcessor) (*types.MsgCreateProcessorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	processor := types.Processor{
		Creator:     msg.Creator,
		ProcessorType: msg.ProcessorType,
		Config:       msg.Config,
		WasmCode:     msg.WasmCode,
	}
	id, err := k.AppendProcessor(ctx, processor)
	if err != nil {
		return nil, err
	}

	return &types.MsgCreateProcessorResponse{Id: id}, nil
}
