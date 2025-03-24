package server

import (
	"context"
	"encoding/json"
	"math/big"

	"cosmossdk.io/math"
	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	dvstypes "github.com/0xPellNetwork/pellapp-sdk/pelldvs/types"
	sdktypes "github.com/0xPellNetwork/pellapp-sdk/types"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"

	"intellix/dvs"
	"intellix/dvs/price/types"
	taskgateway "intellix/gateway"
	pricetypes "intellix/x/price/types"
)

func (d Server) DVSResponsHandler(ctx context.Context, in *types.RequestPriceFeedIn) (*types.ResponsePriceFeedOut, error) {
	pkgCtx := sdktypes.UnwrapContext(ctx)
	//js, _ := json.Marshal(in)
	//d.logger.Info("DvsPostProcessServer.PostProcessRequestPriceFeed called", "data", string(js))

	validatedData, err := pelldvs.GetDvsRequestValidatedData(pkgCtx)
	if err != nil {
		return nil, err
	}

	taskResp, err := dvs.AbiDecodeResponseTaskParam(validatedData.Data)
	if err != nil {
		return nil, err
	}

	// send VoteFinalizedRequestPrice Tx
	_, err = d.sendVoteFinalizedRequestPriceTx(pkgCtx, in, validatedData, taskResp)
	if err != nil {
		return nil, err
	}

	// send gateway
	err = d.sendResponseToGateway(pkgCtx, in, validatedData, taskResp)
	if err != nil {
		return nil, err
	}

	return &types.ResponsePriceFeedOut{}, nil
}

func (d Server) sendVoteFinalizedRequestPriceTx(ctx sdktypes.Context, raw *types.RequestPriceFeedIn, validatedData *dvstypes.RequestPostRequestValidatedData, priceData *contractdataoracle.IDataOracleTaskResponse) (*pricetypes.MsgVoteFinalizedRequestPrice, error) {
	addr, err := d.SenderAddress()
	if err != nil {
		return nil, err
	}

	msg := &pricetypes.MsgVoteFinalizedRequestPrice{
		Sender:                    addr.String(),
		TaskIndex:                 raw.Task.TaskIndex,
		RequestId:                 raw.Task.RequestId,
		FeeToken:                  raw.Task.FeeToken,
		Payment:                   raw.Task.Payment,
		RequestData:               raw.Task.RequestData,
		CallbackAddress:           raw.Task.CallbackAddress,
		CallbackFunctionId:        raw.Task.CallbackFunctionId,
		TaskCreatedBlock:          raw.Task.TaskCreatedBlock,
		QuorumNumbers:             raw.Task.QuorumNumbers,
		QuorumThresholdPercentage: raw.Task.QuorumThresholdPercentage,
		ReferenceTaskIndex:        priceData.ReferenceTaskIndex,
		Price:                     math.NewIntFromBigInt(new(big.Int).SetBytes(priceData.Data)),
	}

	if err := d.SignAndBroadcastTx(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (d Server) sendResponseToGateway(ctx sdktypes.Context, raw *types.RequestPriceFeedIn, validatedData *dvstypes.RequestPostRequestValidatedData, priceData *contractdataoracle.IDataOracleTaskResponse) error {

	nonSignerStakeIndices := make([][]uint32, len(validatedData.NonSignerStakeIndices))
	for i, indices := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices[i] = indices.NonSignerStakeIndice
	}

	req := &taskgateway.RPCVoteFinalizedRequestIn{
		ChainID: ctx.ChainID(),
		TaskRaw: &taskgateway.RPCTaskRaw{
			TaskType:                  taskgateway.TaskTypePriceFeed,
			TaskIndex:                 raw.Task.TaskIndex,
			RequestID:                 raw.Task.RequestId,
			FeeToken:                  raw.Task.FeeToken,
			Payment:                   raw.Task.Payment.String(),
			RequestData:               raw.Task.RequestData,
			CallbackAddress:           raw.Task.CallbackAddress,
			CallbackFunctionID:        raw.Task.CallbackFunctionId,
			TaskCreatedBlock:          raw.Task.TaskCreatedBlock,
			QuorumNumbers:             raw.Task.QuorumNumbers,
			QuorumThresholdPercentage: raw.Task.QuorumThresholdPercentage,
			AdvanceDecode:             raw.Task.AdvanceDecode,
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
		RespToTaskData: priceData.Data,
	}

	reqJs, _ := json.Marshal(req)
	d.logger.Info("DvsPostProcessServer.sendResponseToGateway", "req", string(reqJs))

	return d.taskGatewayClient.RespondToTask(req)
}
