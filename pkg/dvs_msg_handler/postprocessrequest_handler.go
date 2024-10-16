package dvsservermanager

import (
	"github.com/cosmos/cosmos-sdk/codec"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"google.golang.org/grpc"
	"intellix/pkg/dvs_msg_handler/tx"
)

type PostProcessRequestHandler struct {
	Mgr *MsgRouterMgr
}

func NewPostProcessRequestHandler(cdc codec.Codec) grpc1.Server {
	return &PostProcessRequestHandler{
		Mgr: NewMsgRouterMgr(
			tx.NewDefaultDecoder(cdc),
		),
	}
}

func (p *PostProcessRequestHandler) RegisterService(sd *grpc.ServiceDesc, handler interface{}) {
	RegisterServiceRouter(p.Mgr, sd, handler)
}
