package dvsservermanager

import (
	result "intellix/pkg/dvs_msg_handler/result_handler"
	"intellix/pkg/dvs_msg_handler/tx"
	sdktypes "intellix/sdk/types"

	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/cosmos/gogoproto/proto"
	"google.golang.org/grpc"
)

type ProcessRequestHandler struct {
	Mgr           *MsgRouterMgr
	ResultHandler *result.ResultCustomizedMgr
}

func NewProcessRequestHandler(encoder tx.MsgEncoder, resultHandler *result.ResultCustomizedMgr) grpc1.Server {
	return &ProcessRequestHandler{
		Mgr: NewMsgRouterMgr(
			encoder,
			resultHandler,
		),
		ResultHandler: resultHandler,
	}
}

func (p *ProcessRequestHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}

//func (p *ProcessRequestHandler) InvokeRouterByData(ctx pkgcontext.Context, data []byte) ([]byte, error) {
//	res, err := p.Mgr.HandleByData(ctx, data)
//	if err != nil {
//		return nil, err
//	}
//	return res.Data, nil
//}

func (p *ProcessRequestHandler) InvokeRouterRawByData(ctx sdktypes.Context, data []byte) (*result.Result, error) {
	res, err := p.Mgr.HandleByData(ctx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *ProcessRequestHandler) RegisterResultHandler(msg proto.Message, handler result.ResultCustomizedIFace) {
	p.ResultHandler.RegisterCustomizedFunc(msg, handler)
}
