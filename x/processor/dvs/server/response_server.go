package server

import (
	"context"
	"fmt"
	"github.com/IntelliXLabs/iwasm/api"
	taskgateway "intellix/gateway"
	pkgcontext "intellix/pkg/context"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"
	"intellix/x/processor/dvs/types"
	processortypes "intellix/x/processor/types"
)

type ResponseServer struct {
	Server
	runtime api.RuntimeResult
}

func NewResponseServer() types.DVSResponseServer {
	return &ResponseServer{
		runtime: api.NewRuntime(),
	}
}

var _ types.DVSResponseServer = ResponseServer{}

func (r ResponseServer) ResponseScript(ctx context.Context, in *types.RequestScriptIn) (*types.ResponseScriptOut, error) {
	pkgCtx := pkgcontext.UnwrapContext(ctx)
	validatedData, err := r.getDvsRequestValidatedData(pkgCtx)
	if err != nil {
		return nil, err
	}

	err = r.voteData(pkgCtx, in, validatedData.Data)
	if err != nil {
		return nil, err
	}

	err = r.responseToTask(pkgCtx, in, validatedData.Data, validatedData)
	if err != nil {
		return nil, err
	}

	return &types.ResponseScriptOut{}, nil
}

// TODO: move to common code
func (r ResponseServer) getDvsRequestValidatedData(ctx pkgcontext.Context) (*dvstypes.RequestPostRequestValidatedData, error) {
	reqData, ok := ctx.DvsPostResponseData()
	if !ok {
		return nil, fmt.Errorf("not DvsRequestData found")
	}
	validatedDataMsg, err := dvsservermanager.DecodeMsg(reqData)
	if err != nil {
		return nil, err
	}
	validatedData, ok := validatedDataMsg.(*dvstypes.RequestPostRequestValidatedData)
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &dvstypes.RequestPostRequestValidatedData{}, validatedData)
	}
	return validatedData, nil
}

func (r ResponseServer) voteData(ctx pkgcontext.Context, in *types.RequestScriptIn, data []byte) error {
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

func (r ResponseServer) responseToTask(ctx pkgcontext.Context, in *types.RequestScriptIn, data []byte, validatedData *dvstypes.RequestPostRequestValidatedData) error {
	return r.taskGatewayClient.RespondToTask(&taskgateway.RPCVoteFinalizedRequestIn{
		ChainID: ctx.ChainID(),
		TaskRaw: &taskgateway.RPCTaskRaw{
			TaskType:                  taskgateway.TaskTypeProcessor,
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
		},
		ValidatedData: &taskgateway.RPCValidatedData{},
		ScriptResponse: &taskgateway.RPCScriptProcessorResponse{
			Data: data,
		},
	})
}
