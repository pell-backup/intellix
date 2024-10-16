package dvsservermanager

import (
	"github.com/cosmos/cosmos-sdk/codec"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"google.golang.org/grpc"
	"intellix/pkg/dvs_msg_handler/tx"
)

type ProcessRequestHandler struct {
	Mgr *MsgRouterMgr
}

func NewProcessRequestHandler(cdc codec.Codec) grpc1.Server {
	return &ProcessRequestHandler{
		Mgr: NewMsgRouterMgr(
			tx.NewDefaultDecoder(cdc),
		),
	}
}

func (p *ProcessRequestHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}
