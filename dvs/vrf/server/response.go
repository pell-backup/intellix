package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	dvstypes "github.com/0xPellNetwork/pellapp-sdk/pelldvs/types"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"

	"intellix/dvs"
	"intellix/dvs/vrf/types"
	gateway "intellix/gateway/types"
)

// DVSResponsHandler handles the response from the DVS
func (s *Server) DVSResponsHandler(ctx context.Context, in *types.VRFTaskRequest) (*types.DVSResultResponse, error) {
	s.logger.Info("DVSResponsHandler", "TaskIndex", in.TaskMetadata.TaskIndex, "TaskMetadata", fmt.Sprintf("%+v", in.TaskMetadata))
	pkgCtx := sdktypes.UnwrapContext(ctx)

	validatedData, err := pelldvs.GetDvsRequestValidatedData(pkgCtx)
	if err != nil {
		return nil, err
	}

	// Convert []uint32 to bytes
	groupNumbersBytes := make([]byte, len(in.TaskMetadata.GroupNumbers))
	for i, num := range in.TaskMetadata.GroupNumbers {
		groupNumbersBytes[i] = byte(num)
	}

	s.Logger().Info("DVSResponsHandler", "GroupNumbers", string(groupNumbersBytes), "len(validatedData.Data)", len(validatedData.Data))

	taskResp, err := dvs.AbiDecodeResponseTaskParam(validatedData.Data)
	if err != nil {
		return nil, err
	}

	s.Logger().Info("DVSResponsHandler", "TaskResp", taskResp)

	if err := s.sendResponseToGateway(pkgCtx, in, validatedData, taskResp); err != nil {
		s.logger.Error("DVSResponsHandler", "sendResponseToGateway", err)
		return nil, err
	}

	s.logger.Info("RespondToTask Done")
	return &types.DVSResultResponse{}, nil
}

// sendResponseToGateway sends the response to the gateway
func (d *Server) sendResponseToGateway(ctx sdktypes.Context, raw *types.VRFTaskRequest,
	validatedData *dvstypes.RequestPostRequestValidatedData, taskResp *contractdataoracle.IDataOracleTaskResponse) error {

	nonSignerStakeIndices := make([][]uint32, len(validatedData.NonSignerStakeIndices))
	for i, indices := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices[i] = indices.NonSignerStakeIndice
	}

	req := &gateway.RPCVoteFinalizedRequestIn{
		ChainID: ctx.ChainID(),
		TaskRaw: &gateway.RPCTaskRaw{
			TaskType:                  dvs.TaskTypeVRFRandomNumber,
			TaskIndex:                 raw.TaskMetadata.TaskIndex,
			RequestID:                 raw.TaskMetadata.RequestId,
			FeeToken:                  raw.TaskMetadata.FeeToken,
			Payment:                   raw.TaskMetadata.Payment.String(),
			RequestData:               raw.TaskMetadata.RequestData,
			CallbackAddress:           raw.TaskMetadata.CallbackAddress,
			CallbackFunctionID:        raw.TaskMetadata.CallbackFunctionId,
			TaskCreatedBlock:          raw.TaskMetadata.TaskCreatedBlock,
			QuorumNumbers:             raw.TaskMetadata.GroupNumbers,
			QuorumThresholdPercentage: raw.TaskMetadata.GroupThresholdPercentage,
			AdvanceDecode:             raw.TaskMetadata.AdvanceDecode,
		},
		ValidatedData: &gateway.RPCValidatedData{
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
		RespToTaskData: taskResp.Data,
	}

	reqJs, _ := json.Marshal(req)
	d.logger.Info("DvsPostProcessRequestServer.sendResponseToGateway", "req", string(reqJs))

	return d.taskGatewayClient.RespondToDataOracleTask(req)
}
