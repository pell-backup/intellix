package avsi

import (
	"context"
	"encoding/json"
	pkgcontext "intellix/pkg/context"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"

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
	pkgCtx := pkgcontext.NewContext(ctx, nil)
	pkgCtx = pkgCtx.WithBlockHeight(req.Request.Height)
	pkgCtx = pkgCtx.WithChainID(req.Request.ChainId)

	handlerSrc := dvsservermanager.GetProcessRequestHandlerSrc()
	res, err := handlerSrc.InvokeRouterRawByData(pkgCtx, req.Request.Data)
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
	pkgCtx := pkgcontext.NewContext(ctx, nil)
	pkgCtx = pkgCtx.WithBlockHeight(req.DvsRequest.Height)
	pkgCtx = pkgCtx.WithChainID(req.DvsRequest.ChainId)

	reqJs, _ := json.Marshal(req)
	app.logger.Debug("AVSI.ProcessDVSResponse", "req", string(reqJs))

	handlerSrc := dvsservermanager.GetPostProcessRequestHandlerSrc()
	_, err := handlerSrc.InvokeRouterRawByData(pkgCtx, req.DvsRequest.Data, dvstypes.NewValidatedResponse(req.DvsResponse))
	if err != nil {
		app.logger.Error("post request error", "err", err)
		return nil, err
	}

	return &avsitypes.ResponseProcessDVSResponse{}, nil
}
