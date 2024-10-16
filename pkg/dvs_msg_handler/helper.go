package dvsservermanager

import (
	"github.com/cosmos/cosmos-sdk/codec"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"sync"
)

type dvsMsgHelper struct {
	lock                      sync.RWMutex
	PostProcessRequestHandler grpc1.Server
	ProcessRequestHandler     grpc1.Server
}

var helper = &dvsMsgHelper{
	lock: sync.RWMutex{},
}

func InitDvsMsgHelper(cdc codec.Codec) {
	helper.lock.Lock()
	defer helper.lock.Unlock()

	if helper.ProcessRequestHandler == nil {
		helper.ProcessRequestHandler = NewProcessRequestHandler(cdc)
	}
	if helper.PostProcessRequestHandler == nil {
		helper.PostProcessRequestHandler = NewPostProcessRequestHandler(cdc)
	}
}

func GetPostProcessRequestHandler() grpc1.Server {
	return helper.PostProcessRequestHandler
}

func GetProcessRequestHandler() grpc1.Server {
	return helper.ProcessRequestHandler
}
