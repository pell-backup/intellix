package dvs

import (
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"intellix/dvs/vrf/handler"
	"intellix/dvs/vrf/server"
	"intellix/dvs/vrf/types"
	dvsservermanager "intellix/sdk/dvs_msg_handler"
)

// AppModule implements an application module for the dvs module.
type AppModule struct {
	server         server.Server
	RequestServer  grpc1.Server
	ResponseServer grpc1.Server
}

// NewAppModule creates a new AppModule object
func NewAppModule(s server.Server) AppModule {
	return AppModule{
		server:         s,
		RequestServer:  dvsservermanager.GetProcessRequestHandler(),
		ResponseServer: dvsservermanager.GetPostProcessRequestHandler(),
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices() {
	requestServer := server.NewRequestServer(am.server)

	// register dvs-msg handler server
	types.RegisterVRFMsgServerServer(am.RequestServer, requestServer)

	// register dvs-msg result handler
	if r, ok := am.RequestServer.(*dvsservermanager.ProcessRequestHandler); ok {
		r.RegisterResultHandler(
			&types.GenerateRandomNumberResponse{}, handler.NewVRFResultHandler(),
		)
	}

}
