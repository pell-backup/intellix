package app

import (
	"context"
	pkgcontext "intellix/pkg/context"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"

	avsitypes "github.com/0xPellNetwork/pelldvs/avsi/types"
)

func (p *PellApp) Info(ctx context.Context, info *avsitypes.RequestInfo) (*avsitypes.ResponseInfo, error) {
	return &avsitypes.ResponseInfo{
		Version:         "1.0.0",
		LastBlockHeight: 0,
	}, nil
}

func (p *PellApp) Query(ctx context.Context, query *avsitypes.RequestQuery) (*avsitypes.ResponseQuery, error) {
	return &avsitypes.ResponseQuery{
		Code: avsitypes.CodeTypeOK,
	}, nil
}

func (p *PellApp) ProcessRequest(ctx context.Context, req *avsitypes.RequestProcessRequest) (*avsitypes.ResponseProcessRequest, error) {
	// new SDK context
	pkgCtx := pkgcontext.NewContext(ctx, nil)
	pkgCtx = pkgCtx.WithBlockHeight(req.Request.Height)
	pkgCtx = pkgCtx.WithChainID(req.Request.ChainId)

	handlerSrc := dvsservermanager.GetProcessRequestHandlerSrc()
	res, err := handlerSrc.InvokeRouterRawByData(pkgCtx, req.Request.Data)
	if err != nil {
		p.logger.Error("process request error", "err", err)
		return nil, err
	}

	return &avsitypes.ResponseProcessRequest{
		Response:       res.CustomData,
		ResponseDigest: res.CustomDigest,
	}, err
}

func (p *PellApp) PostRequest(ctx context.Context, req *avsitypes.RequestPostRequest) (*avsitypes.ResponsePostRequest, error) {
	// new SDK context
	pkgCtx := pkgcontext.NewContext(ctx, nil)
	pkgCtx = pkgCtx.WithBlockHeight(req.DvsRequest.Height)
	pkgCtx = pkgCtx.WithChainID(req.DvsRequest.ChainId)

	handlerSrc := dvsservermanager.GetPostProcessRequestHandlerSrc()
	_, err := handlerSrc.InvokeRouterRawByData(pkgCtx, req.DvsRequest.Data, dvstypes.NewValidatedResponse(req.ValidatedResponse))
	if err != nil {
		p.logger.Error("post request error", "err", err)
		return nil, err
	}

	return &avsitypes.ResponsePostRequest{}, nil
}
