package dvsservermanager

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"google.golang.org/grpc"
	"intellix/pkg/dvs_msg_handler/tx"
	"strings"
)

type PostProcessRequestHandler struct {
	Mgr *MsgRouterMgr
}

func NewPostProcessRequestHandler(encoder tx.MsgEncoder) grpc1.Server {
	return &PostProcessRequestHandler{
		Mgr: NewMsgRouterMgr(encoder, func(msg sdk.Msg) string {
			// router postProcessRequestReq by processRequestReq
			r := sdk.MsgTypeURL(msg)
			if strings.HasPrefix(r, "ProcessRequest") {
				return strings.ReplaceAll(r, "ProcessRequest", "PostProcessRequest")
			}
			return r
		}),
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
