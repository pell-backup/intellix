package dvsservermanager

import (
	pkgcontext "intellix/pkg/context"
	result "intellix/pkg/dvs_msg_handler/result_handler"
	"intellix/pkg/dvs_msg_handler/tx"

	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
)

type PostProcessRequestHandler struct {
	Mgr           *MsgRouterMgr
	ResultHandler *result.ResultCustomizedMgr
}

func NewPostProcessRequestHandler(encoder tx.MsgEncoder, resultHandler *result.ResultCustomizedMgr) grpc1.Server {
	return &PostProcessRequestHandler{
		Mgr: NewMsgRouterMgr(
			encoder,
			resultHandler,
		),
		ResultHandler: resultHandler,
	}
}

func (p *PostProcessRequestHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}

// InvokeRouterRawByData
// requestData: binary data from processRequestData, for found router and dispatcher
// reqMsg: post-process-response data, attached to context
func (p *PostProcessRequestHandler) InvokeRouterRawByData(ctx pkgcontext.Context, requestData []byte, postProcessResponseMsg sdk.Msg) (*result.Result, error) {
	postResponseData, err := p.Mgr.encoder.EncodeMsgs(postProcessResponseMsg)
	if err != nil {
		return nil, err
	}
	ctx = ctx.WithDvsPostResponseData(postResponseData)

	return p.Mgr.HandleByData(ctx, requestData)
}

func (p *PostProcessRequestHandler) RegisterResultHandler(msg proto.Message, handler result.ResultCustomizedIFace) {
	p.ResultHandler.RegisterCustomizedFunc(msg, handler)
}
