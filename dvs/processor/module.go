package dvs

import (
	grpc1 "github.com/cosmos/gogoproto/grpc"

	resulthandlers "intellix/dvs/processor/result_handlers"
	"intellix/dvs/processor/server"
	"intellix/dvs/processor/types"
	dvsservermanager "intellix/sdk/dvs_msg_handler"
)

type AppModule struct {
	server         server.Server
	RequestServer  grpc1.Server
	ResponseServer grpc1.Server
}

func NewAppModule(s server.Server) *AppModule {
	return &AppModule{
		server:         s,
		RequestServer:  dvsservermanager.GetProcessRequestHandler(),
		ResponseServer: dvsservermanager.GetPostProcessRequestHandler(),
	}
}

func (am AppModule) RegisterServices() {
	reqServer := server.NewRequestServer(am.server)
	respServer := server.NewResponseServer(am.server)

	types.RegisterDVSRequestServer(am.RequestServer, reqServer)
	types.RegisterDVSResponseServer(am.ResponseServer, respServer)

	if r, ok := am.RequestServer.(*dvsservermanager.ProcessRequestHandler); ok {
		r.RegisterResultHandler(
			&types.RequestScriptOut{}, resulthandlers.NewProcessorRequestResHandler(),
		)
	}
}
