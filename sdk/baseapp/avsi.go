package baseapp

import (
	"context"
	"encoding/json"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"
	sdktypes "intellix/sdk/types"

	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
)

func (app *BaseApp) Info(ctx context.Context, info *avsitypes.RequestInfo) (*avsitypes.ResponseInfo, error) {
	return &avsitypes.ResponseInfo{
		Version:         "1.0.0",
		LastBlockHeight: 0,
	}, nil
}

func (app *BaseApp) Query(ctx context.Context, query *avsitypes.RequestQuery) (*avsitypes.ResponseQuery, error) {
	return &avsitypes.ResponseQuery{
		Code: avsitypes.CodeTypeOK,
	}, nil
}

func (app *BaseApp) ProcessDVSRequest(ctx context.Context, req *avsitypes.RequestProcessDVSRequest) (*avsitypes.ResponseProcessDVSRequest, error) {
	// new SDK context
	sdkCtx := sdktypes.NewContext(ctx)
	sdkCtx = sdkCtx.WithChainID(req.Request.ChainId).
		WithHeight(req.Request.Height).
		WithGroupNumbers(req.Request.GroupNumbers).
		WithGroupThresholdPercentages(req.Request.GroupThresholdPercentages).
		WithOperator(req.Operator)

	reqJs, _ := json.Marshal(req)
	app.logger.Debug("AVSI.ProcessDVSRequest", "req", string(reqJs))

	handlerSrc := dvsservermanager.GetProcessRequestHandlerSrc()
	res, err := handlerSrc.InvokeRouterRawByData(sdkCtx, req.Request.Data)
	if err != nil {
		app.logger.Error("process request error", "err", err)
		return nil, err
	}

	return &avsitypes.ResponseProcessDVSRequest{
		Response:       res.CustomData,
		ResponseDigest: res.CustomDigest,
	}, err
}

func (app *BaseApp) ProcessDVSResponse(ctx context.Context, req *avsitypes.RequestProcessDVSResponse) (*avsitypes.ResponseProcessDVSResponse, error) {
	// new SDK context
	sdkCtx := sdktypes.NewContext(ctx)
	sdkCtx = sdkCtx.WithChainID(req.DvsRequest.ChainId).
		WithHeight(req.DvsRequest.Height).
		WithGroupNumbers(req.DvsRequest.GroupNumbers).
		WithGroupThresholdPercentages(req.DvsRequest.GroupThresholdPercentages)

	reqJs, _ := json.Marshal(req)
	app.logger.Debug("AVSI.ProcessDVSResponse", "req", string(reqJs))

	handlerSrc := dvsservermanager.GetPostProcessRequestHandlerSrc()
	_, err := handlerSrc.InvokeRouterRawByData(sdkCtx, req.DvsRequest.Data, dvstypes.NewValidatedResponse(req.DvsResponse))
	if err != nil {
		app.logger.Error("post request error", "err", err)
		return nil, err
	}

	return &avsitypes.ResponseProcessDVSResponse{}, nil
}
