package server

import (
	"context"
	"intellix/dvs"
	taskgateway "intellix/gateway/types"

	"intellix/dvs/processor/types"
	"intellix/sdk/pelldvs"
	dvstypes "intellix/sdk/pelldvs/types"
	sdktypes "intellix/sdk/types"
	"intellix/sdk/utils"
	processortypes "intellix/x/processor/types"
)

type ResponseServer struct {
	Server
}

func NewResponseServer(s Server) types.DVSResponseServer {
	return &ResponseServer{
		Server: s,
	}
}

var _ types.DVSResponseServer = ResponseServer{}

func (r ResponseServer) ResponseScript(ctx context.Context, in *types.RequestScriptIn) (*types.ResponseScriptOut, error) {
	pkgCtx := sdktypes.UnwrapContext(ctx)
	validatedData, err := pelldvs.GetDvsRequestValidatedData(pkgCtx)
	if err != nil {
		return nil, err
	}

	// decode data from abi-encoded-data
	scriptOutData, err := utils.AbiDecodeResponseTaskParam(validatedData.Data)
	if err != nil {
		return nil, err
	}

	err = r.voteData(pkgCtx, in, scriptOutData.Data)
	if err != nil {
		return nil, err
	}

	err = r.responseToTask(pkgCtx, in, scriptOutData.Data, validatedData)
	if err != nil {
		return nil, err
	}

	return &types.ResponseScriptOut{}, nil
}

func (r ResponseServer) voteData(ctx sdktypes.Context, in *types.RequestScriptIn, data []byte) error {
	msg := &processortypes.MsgVoteResponseProcessor{
		TaskIndex:                 in.TaskIndex,
		RequestId:                 in.RequestId,
		FeeToken:                  in.FeeToken,
		Payment:                   in.Payment,
		RequestData:               in.RequestData,
		CallbackAddress:           in.CallbackAddress,
		CallbackFunctionId:        in.CallbackFunctionId,
		TaskCreatedBlock:          in.TaskCreatedBlock,
		QuorumNumbers:             in.QuorumNumbers,
		QuorumThresholdPercentage: in.QuorumThresholdPercentage,
		ScriptId:                  in.ScriptId,
		ScriptAggredResp:          data,
	}
	if err := r.SignAndBroadcastTx(ctx, msg); err != nil {
		return err
	}
	return nil
}

func (r ResponseServer) responseToTask(ctx sdktypes.Context, in *types.RequestScriptIn, data []byte, validatedData *dvstypes.RequestPostRequestValidatedData) error {
	var nonSignerStakeIndices [][]uint32
	for _, v := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices = append(nonSignerStakeIndices, v.NonSignerStakeIndice)
	}

	return r.taskGatewayClient.RespondToPriceTask(&taskgateway.RPCVoteFinalizedRequestIn{
		ChainID: ctx.ChainID(),
		TaskRaw: &taskgateway.RPCTaskRaw{
			TaskType:                  dvs.TaskTypeProcessor,
			TaskIndex:                 in.TaskIndex,
			RequestID:                 in.RequestId,
			FeeToken:                  in.FeeToken,
			Payment:                   in.Payment.String(),
			RequestData:               in.RequestData,
			CallbackAddress:           in.CallbackAddress,
			CallbackFunctionID:        in.CallbackFunctionId,
			TaskCreatedBlock:          in.TaskCreatedBlock,
			QuorumNumbers:             in.QuorumNumbers,
			QuorumThresholdPercentage: in.QuorumThresholdPercentage,
			AdvanceDecode:             in.AdvanceDecode,
		},
		ValidatedData: &taskgateway.RPCValidatedData{
			Data:                         validatedData.Data,
			Error:                        validatedData.Error,
			Hash:                         validatedData.Hash,
			NonSignersPubkeysG1:          validatedData.NonSignersPubkeysG1,
			QuorumApksG1:                 validatedData.QuorumApksG1,
			SignersApkG2:                 validatedData.SignersApkG2,
			SignersAggSigG1:              validatedData.SignersAggSigG1,
			NonSignerQuorumBitmapIndices: validatedData.NonSignerQuorumBitmapIndices,
			QuorumApkIndices:             validatedData.QuorumApkIndices,
			TotalStakeIndices:            validatedData.TotalStakeIndices,
			NonSignerStakeIndices:        nonSignerStakeIndices,
		},
		RespToTaskData: data,
	})
}
