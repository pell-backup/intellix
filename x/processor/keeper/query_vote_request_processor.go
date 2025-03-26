package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"intellix/x/processor/types"
)

// QueryVoteRequestProcessor queries vote request processor
func (k Keeper) QueryVoteRequestProcessor(goCtx context.Context, req *types.QueryVoteRequestProcessorRequest) (*types.QueryVoteRequestProcessorResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(goCtx)

	if req.RequestId == nil || req.OperatorId == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	voteRequestProcessor, found := k.GetVoteRequestProcessor(sdkCtx, req.RequestId, req.OperatorId)
	if !found {
		return nil, status.Error(codes.NotFound, "vote request processor not found")
	}

	return &types.QueryVoteRequestProcessorResponse{
		TaskIndex:                 voteRequestProcessor.TaskIndex,
		RequestId:                 voteRequestProcessor.RequestId,
		FeeToken:                  voteRequestProcessor.FeeToken,
		Payment:                   voteRequestProcessor.Payment,
		RequestData:               voteRequestProcessor.RequestData,
		CallbackAddress:           voteRequestProcessor.CallbackAddress,
		CallbackFunctionId:        voteRequestProcessor.CallbackFunctionId,
		TaskCreatedBlock:          voteRequestProcessor.TaskCreatedBlock,
		QuorumNumbers:             voteRequestProcessor.QuorumNumbers,
		QuorumThresholdPercentage: voteRequestProcessor.QuorumThresholdPercentage,
		ScriptId:                  voteRequestProcessor.ScriptId,
		ScriptResp:                voteRequestProcessor.ScriptResp,
		Sender:                    voteRequestProcessor.Sender,
		OperatorId:                voteRequestProcessor.OperatorId,
		BlsSignature:              voteRequestProcessor.BlsSignature,
	}, nil
}

// GetVoteRequestProcessor returns vote request processor
func (k Keeper) GetVoteRequestProcessor(ctx sdk.Context, requestId []byte, operatorId string) (*types.MsgVoteRequestProcessor, bool) {
	store := k.storeService.OpenKVStore(ctx)

	data, err := store.Get(types.MsgVoteRequestProcessorKey(requestId, operatorId))
	if err != nil {
		return nil, false
	}

	var msg types.MsgVoteRequestProcessor
	if err := k.cdc.Unmarshal(data, &msg); err != nil {
		return nil, false
	}

	return &msg, true
}
