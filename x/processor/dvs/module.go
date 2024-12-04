package dvs

import (
	grpc1 "github.com/cosmos/gogoproto/grpc"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	resulthandlers "intellix/x/price/dvs/result_handlers"
	"intellix/x/processor/dvs/server"
	"intellix/x/processor/dvs/types"
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
	reqServer := server.NewRequestServer()
	respServer := server.NewResponseServer()

	types.RegisterDVSRequestServer(am.RequestServer, reqServer)
	types.RegisterDVSResponseServer(am.ResponseServer, respServer)

	if r, ok := am.RequestServer.(*dvsservermanager.ProcessRequestHandler); ok {
		r.RegisterResultHandler(
			&types.RequestScriptOut{}, resulthandlers.NewProcessorRequestResHandler(),
		)
	}
}
