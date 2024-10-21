package dvsservermanager

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
	result "intellix/pkg/dvs_msg_handler/result_handler"
	"intellix/pkg/dvs_msg_handler/tx"
	"strings"
)

type PostProcessRequestHandler struct {
	Mgr           *MsgRouterMgr
	ResultHandler *result.ResultCustomizedMgr
}

func NewPostProcessRequestHandler(encoder tx.MsgEncoder, resultHandler *result.ResultCustomizedMgr) grpc1.Server {
	return &PostProcessRequestHandler{
		Mgr: NewMsgRouterMgr(
			encoder,
			func(msg sdk.Msg) string {
				// router postProcessRequestReq by processRequestReq
				r := sdk.MsgTypeURL(msg)
				if strings.HasPrefix(r, "ProcessRequest") {
					return strings.ReplaceAll(r, "ProcessRequest", "PostProcessRequest")
				}
				return r
			},
			resultHandler,
		),
		ResultHandler: resultHandler,
	}
}

func (p *PostProcessRequestHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}

func (p *PostProcessRequestHandler) InvokeRouterByData(sdkCtx sdk.Context, sourceData []byte, reqMsg sdk.Msg) ([]byte, error) {
	handler := p.Mgr.GetHandlerByData(sourceData)
	if handler == nil {
		return nil, fmt.Errorf("no handler found for data: %s", string(sourceData))
	}

	res, err := handler(sdkCtx, reqMsg)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (p *PostProcessRequestHandler) InvokeRouterRawByData(sdkCtx sdk.Context, sourceData []byte, reqMsg sdk.Msg) (*result.Result, error) {
	handler := p.Mgr.GetHandlerByData(sourceData)
	if handler == nil {
		return nil, fmt.Errorf("no handler found for data: %s", string(sourceData))
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
