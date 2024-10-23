package dvsservermanager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
	result "intellix/pkg/dvs_msg_handler/result_handler"
	"intellix/pkg/dvs_msg_handler/tx"
)

type PostProcessRequestHandler struct {
	Mgr           *MsgRouterMgr
	ResultHandler *result.ResultCustomizedMgr
}

func NewPostProcessRequestHandler(encoder tx.MsgEncoder, resultHandler *result.ResultCustomizedMgr) grpc1.Server {
	return &PostProcessRequestHandler{
		Mgr: NewMsgRouterMgr(
			encoder,
			nil,
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
func (p *PostProcessRequestHandler) InvokeRouterRawByData(sdkCtx sdk.Context, requestData []byte, postProcessResponseMsg sdk.Msg) (*result.Result, error) {
	postResponseData, err := p.Mgr.encoder.EncodeMsgs(postProcessResponseMsg)
	if err != nil {
		return nil, err
	}
	sdkCtx = CtxWithDvsPostResponseData(sdkCtx, postResponseData)

	return p.Mgr.HandleByData(sdkCtx, requestData)
}

func (p *PostProcessRequestHandler) RegisterResultHandler(msg proto.Message, handler result.ResultCustomizedIFace) {
	p.ResultHandler.RegisterCustomizedFunc(msg, handler)
}
