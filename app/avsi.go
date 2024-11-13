package app

import (
	"context"
	pkgcontext "intellix/pkg/context"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/pkg/pelldvs/types"

	avsi "github.com/0xPellNetwork/pelldvs/application"
)

func (app *App) ProcessRequest(ctx context.Context, req *avsi.RequestProcessRequest) (*avsi.ResponseProcessRequest, error) {
	// new SDK context
	pkgCtx := pkgcontext.NewContext(ctx, nil)
	pkgCtx = pkgCtx.WithBlockHeight(req.Request.Height)
	pkgCtx = pkgCtx.WithChainID(req.Request.ChainID.String())

	handlerSrc := dvsservermanager.GetProcessRequestHandlerSrc()
	res, err := handlerSrc.InvokeRouterRawByData(pkgCtx, req.Request.Data)
	if err != nil {
		return nil, err
	}

	return &avsi.ResponseProcessRequest{
		Reponse:        res.CustomData,
		ResponseDigest: res.CustomDigest,
	}, err
}

func (app *App) PostRequest(ctx context.Context, req *avsi.RequestPostRequest) (*avsi.ResponsePostRequest, error) {
	// new SDK context
	pkgCtx := pkgcontext.NewContext(ctx, nil)
	pkgCtx = pkgCtx.WithBlockHeight(req.Request.Height)
	pkgCtx = pkgCtx.WithChainID(req.Request.ChainID.String())

	handlerSrc := dvsservermanager.GetPostProcessRequestHandlerSrc()
	res, err := handlerSrc.InvokeRouterRawByData(pkgCtx, req.Request.Data, dvstypes.NewValidatedResponse(&req.Response))
	if err != nil {
		return nil, err
	}

	return &avsi.ResponsePostRequest{
		Receipt: res.CustomData,
	}, nil
}
