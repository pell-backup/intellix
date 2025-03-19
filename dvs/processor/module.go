package dvs

import (
	sdkservice "github.com/0xPellNetwork/pellapp-sdk/service"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	resulthandlers "intellix/dvs/processor/result_handlers"
	"intellix/dvs/processor/server"
	"intellix/dvs/processor/types"
)

type AppModule struct {
	server server.Server
}

func NewAppModule(s server.Server) *AppModule {
	return &AppModule{
		server: s,
	}
}

func (am AppModule) RegisterServices(router *sdkservice.MsgRouter) {
	configurator := router.GetConfigurator()
	// register dvs-msg handler server
	types.RegisterDVSRequestServer(configurator, am.server)

	// register dvs-msg result handler
	configurator.RegisterResultMsgExtractor(
		&types.RequestScriptOut{}, resulthandlers.NewProcessorRequestResHandler(),
	)
}

func (am AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}
