package dvs

import (
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	resulthandlers "intellix/x/price/dvs/result_handlers"
	"intellix/x/price/dvs/server"
	"intellix/x/price/dvs/types"

	grpc1 "github.com/cosmos/gogoproto/grpc"
)

// AppModule implements an application module for the dvs module.
type AppModule struct {
	server                   server.Server
	ProcessRequestServer     grpc1.Server
	PostProcessRequestServer grpc1.Server
}

// NewAppModule creates a new AppModule object
func NewAppModule(s server.Server) AppModule {
	return AppModule{
		server:                   s,
		ProcessRequestServer:     dvsservermanager.GetProcessRequestHandler(),
		PostProcessRequestServer: dvsservermanager.GetPostProcessRequestHandler(),
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices() {
	dvsProcessRequestServer := server.NewRequestServer(am.server)
	dvsPostProcessRequestServer := server.NewResponseServer(am.server)

	// register dvs-msg handler server
	types.RegisterDvsProcessRequestServer(am.ProcessRequestServer, dvsProcessRequestServer)
	types.RegisterDvsPostProcessRequestServer(am.PostProcessRequestServer, dvsPostProcessRequestServer)

	// register dvs-msg result handler
	if r, ok := am.ProcessRequestServer.(*dvsservermanager.ProcessRequestHandler); ok {
		r.RegisterResultHandler(
			&types.ProcessRequestPriceFeedOut{}, resulthandlers.NewProcessRequestPriceFeedResultHandler(),
		)
	}

}
