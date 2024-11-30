package server

import (
	"context"
	"encoding/json"
	"fmt"
	taskgateway "intellix/gateway"
	pkgcontext "intellix/pkg/context"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"
	"intellix/x/price/dvs/types"
	pricetypes "intellix/x/price/types"
	"math/big"

	"cosmossdk.io/math"
	contractDataOracle "github.com/IntelliXLabs/price-oracle-dvs/bindings/DataOracle"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func (d ResponseServer) PostProcessRequestPriceFeed(ctx context.Context, in *types.ProcessRequestPriceFeedIn) (*types.PostProcessRequestPriceFeedOut, error) {
	pkgCtx := pkgcontext.UnwrapContext(ctx)
	js, _ := json.Marshal(in)
	d.logger.Info("DvsPostProcessRequestServer.PostProcessRequestPriceFeed called", "data", string(js))

	validatedData, err := d.getDvsRequestValidatedData(pkgCtx)
	if err != nil {
		return nil, err
	}
	taskResp, err := d.decodePackedPriceFeedData(validatedData.Data)
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

	return &types.PostProcessRequestPriceFeedOut{}, nil
}

func (d ResponseServer) sendVoteFinalizedRequestPriceTx(ctx pkgcontext.Context, raw *types.ProcessRequestPriceFeedIn, validatedData *dvstypes.RequestPostRequestValidatedData, priceData *contractDataOracle.IDataOracleTaskResponse) (*pricetypes.MsgVoteFinalizedRequestPrice, error) {
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

func (d ResponseServer) sendResponseToGateway(ctx pkgcontext.Context, raw *types.ProcessRequestPriceFeedIn, validatedData *dvstypes.RequestPostRequestValidatedData, priceData *contractDataOracle.IDataOracleTaskResponse) error {

	nonSignerStakeIndices := make([][]uint32, len(validatedData.NonSignerStakeIndices))
	for i, indices := range validatedData.NonSignerStakeIndices {
		nonSignerStakeIndices[i] = indices.NonSignerStakeIndice
	}

	req := &taskgateway.RPCVoteFinalizedRequestPrice{
		ChainID: ctx.ChainID(),
		TaskRaw: &taskgateway.RPCTaskRaw{
			TaskType:                  types.TaskTypePriceFeed,
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
		PriceFeedResponse: &taskgateway.RPCPriceFeedResponse{
			ReferenceTaskIndex: priceData.ReferenceTaskIndex,
			Price:              new(big.Int).SetBytes(priceData.Data).String(),
		},
	}

	reqJs, _ := json.Marshal(req)
	d.logger.Info("DvsPostProcessRequestServer.sendResponseToGateway", "req", string(reqJs))

	return d.Server.taskGatewayClient.RespondToTask(req)
}

func (d ResponseServer) decodePackedPriceFeedData(data []byte) (*contractDataOracle.IDataOracleTaskResponse, error) {
	taskResponseType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "referenceTaskIndex",
			Type: "uint32",
		},
		{
			Name: "data",
			Type: "bytes",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ABI type: %w", err)
	}

	arguments := abi.Arguments{
		{
			Type: taskResponseType,
		},
	}

	// decode
	values, err := arguments.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack data: %w", err)
	}

	if len(values) != 1 {
		return nil, fmt.Errorf("unexpected number of values: got %d, want 1", len(values))
	}

	r, ok := values[0].(struct {
		ReferenceTaskIndex uint32 `json:"referenceTaskIndex"`
		Data               []byte `json:"data"`
	})
	d.logger.Debug("decodePackedPriceFeedData", "r", fmt.Sprintf("%+v", r))
	if !ok {
		return nil, fmt.Errorf("expected %T, got %T", &contractDataOracle.IDataOracleTaskResponse{}, values[0])
	}

	return &contractDataOracle.IDataOracleTaskResponse{
		ReferenceTaskIndex: r.ReferenceTaskIndex,
		Data:               r.Data,
	}, nil
}

func (d ResponseServer) getDvsRequestValidatedData(ctx pkgcontext.Context) (*dvstypes.RequestPostRequestValidatedData, error) {
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
