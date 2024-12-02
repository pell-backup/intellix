package dvs

import (
	grpc1 "github.com/cosmos/gogoproto/grpc"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	"intellix/x/processor/dvs/server"
	"intellix/x/processor/dvs/types"
)

type AppModule struct {
	RequestServer grpc1.Server
}

func NewAppModule() *AppModule {
	return &AppModule{
		RequestServer: dvsservermanager.GetProcessRequestHandler(),
	}
}

func (am AppModule) RegisterServices() {
	reqServer := server.NewRequestServer()

	types.RegisterDVSRequestServer(am.RequestServer, reqServer)
}
