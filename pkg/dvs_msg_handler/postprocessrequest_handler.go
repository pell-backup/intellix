package dvsservermanager

import (
	"fmt"
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

func (p *PostProcessRequestHandler) InvokeRouterRawByData(sdkCtx sdk.Context, reqMsg sdk.Msg) (*result.Result, error) {
	data, err := p.Mgr.encoder.EncodeMsgs(reqMsg)
	if err != nil {
		return nil, err
	}

	handler := p.Mgr.GetHandlerByData(data)
	if handler == nil {
		return nil, fmt.Errorf("no handler found for data: %s", string(data))
	}

	res, err := handler(sdkCtx, reqMsg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (p *PostProcessRequestHandler) RegisterResultHandler(msg proto.Message, handler result.ResultCustomizedIFace) {
	p.ResultHandler.RegisterCustomizedFunc(msg, handler)
}
