package dvsservermanager

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"intellix/pkg/dvs_msg_handler/tx"
	"sync"
)

type dvsMsgHelper struct {
	lock sync.RWMutex

	cdc     codec.Codec
	encoder tx.MsgEncoder

	PostProcessRequestHandler grpc1.Server
	ProcessRequestHandler     grpc1.Server
}

// TODO: use dependency auto inject
var helper = &dvsMsgHelper{
	lock: sync.RWMutex{},
}

func InitDvsMsgHelper(cdc codec.Codec) {
	helper.lock.Lock()
	defer helper.lock.Unlock()

	if helper.cdc == nil {
		helper.cdc = cdc
	}
	if helper.encoder == nil {
		helper.encoder = tx.NewDefaultDecoder(helper.cdc)
	}

	if helper.ProcessRequestHandler == nil {
		helper.ProcessRequestHandler = NewProcessRequestHandler(helper.encoder)
	}
	if helper.PostProcessRequestHandler == nil {
		helper.PostProcessRequestHandler = NewPostProcessRequestHandler(helper.encoder)
	}
}

func GetPostProcessRequestHandler() grpc1.Server {
	return helper.PostProcessRequestHandler
}

func GetProcessRequestHandler() grpc1.Server {
	return helper.ProcessRequestHandler
}

func GetPostProcessRequestHandlerSrc() *PostProcessRequestHandler {
	return helper.PostProcessRequestHandler.(*PostProcessRequestHandler)
}

func GetProcessRequestHandlerSrc() *ProcessRequestHandler {
	return helper.ProcessRequestHandler.(*ProcessRequestHandler)
}

func EncodeMsgs(msgs ...sdk.Msg) ([]byte, error) {
	return helper.encoder.EncodeMsgs(msgs...)
}
