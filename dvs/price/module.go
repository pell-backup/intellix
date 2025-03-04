package dvs

import (
	grpc1 "github.com/cosmos/gogoproto/grpc"

	resulthandlers "intellix/dvs/price/result_handlers"
	"intellix/dvs/price/server"
	"intellix/dvs/price/types"
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
	dvsProcessRequestServer := server.NewRequestServer(am.server)
	dvsPostProcessRequestServer := server.NewResponseServer(am.server)

	// register dvs-msg handler server
	types.RegisterDVSRequestServer(am.RequestServer, dvsProcessRequestServer)
	types.RegisterDVSResponseServer(am.ResponseServer, dvsPostProcessRequestServer)

	// register dvs-msg result handler
	if r, ok := am.RequestServer.(*dvsservermanager.ProcessRequestHandler); ok {
		r.RegisterResultHandler(
			&types.RequestPriceFeedOut{}, resulthandlers.NewProcessRequestPriceFeedResultHandler(),
		)
	}

}
