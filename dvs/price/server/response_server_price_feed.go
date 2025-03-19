package server

import (
	"context"
	"encoding/json"
	"intellix/dvs"
	gateway "intellix/gateway/types"
	"math/big"

	"cosmossdk.io/math"
	contractdataoracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"

	"intellix/dvs/price/types"
	"intellix/sdk/pelldvs"
	dvstypes "intellix/sdk/pelldvs/types"
	sdktypes "intellix/sdk/types"
	"intellix/sdk/utils"
	pricetypes "intellix/x/price/types"
)

func (d ResponseServer) ResponsePriceFeed(ctx context.Context, in *types.RequestPriceFeedIn) (*types.ResponsePriceFeedOut, error) {
	pkgCtx := sdktypes.UnwrapContext(ctx)
	//js, _ := json.Marshal(in)
	//d.logger.Info("DvsPostProcessRequestServer.PostProcessRequestPriceFeed called", "data", string(js))

	validatedData, err := pelldvs.GetDvsRequestValidatedData(pkgCtx)
	if err != nil {
		return nil, err
	}

	taskResp, err := utils.AbiDecodeResponseTaskParam(validatedData.Data)
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

func (d ResponseServer) sendVoteFinalizedRequestPriceTx(ctx sdktypes.Context, raw *types.RequestPriceFeedIn, validatedData *dvstypes.RequestPostRequestValidatedData, priceData *contractdataoracle.IDataOracleTaskResponse) (*pricetypes.MsgVoteFinalizedRequestPrice, error) {
	addr, err := d.Server.SenderAddress()
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

	if err := d.Server.SignAndBroadcastTx(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (d ResponseServer) sendResponseToGateway(ctx sdktypes.Context, raw *types.RequestPriceFeedIn, validatedData *dvstypes.RequestPostRequestValidatedData, priceData *contractdataoracle.IDataOracleTaskResponse) error {

	nonSignerStakeIndices := make([][]uint32, len(validatedData.NonSignerStakeIndices))
	for i, indices := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices[i] = indices.NonSignerStakeIndice
	}

	req := &gateway.RPCVoteFinalizedRequestIn{
		ChainID: ctx.ChainID(),
		TaskRaw: &gateway.RPCTaskRaw{
			TaskType:                  dvs.TaskTypePriceFeed,
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
		RespToTaskData: priceData.Data,
	}

	reqJs, _ := json.Marshal(req)
	d.logger.Info("DvsPostProcessRequestServer.sendResponseToGateway", "req", string(reqJs))

	return d.Server.taskGatewayClient.RespondToPriceTask(req)
}
