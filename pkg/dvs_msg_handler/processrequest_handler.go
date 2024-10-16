package dvsservermanager

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"google.golang.org/grpc"
	"intellix/pkg/dvs_msg_handler/tx"
)

type ProcessRequestHandler struct {
	Mgr *MsgRouterMgr
}

func NewProcessRequestHandler(encoder tx.MsgEncoder) grpc1.Server {
	return &ProcessRequestHandler{
		Mgr: NewMsgRouterMgr(encoder, nil),
	}
}

func (p *ProcessRequestHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}

func (p *ProcessRequestHandler) InvokeRouterByData(sdkCtx sdk.Context, data []byte) ([]byte, error) {
	res, err := p.Mgr.HandleByData(sdkCtx, data)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}
