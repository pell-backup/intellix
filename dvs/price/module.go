package dvs

import (
	sdkservice "github.com/0xPellNetwork/pellapp-sdk/service"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	resulthandlers "intellix/dvs/price/result_handlers"
	"intellix/dvs/price/server"
	"intellix/dvs/price/types"
)

// AppModule implements an application module for the dvs module.
type AppModule struct {
	server server.Server
}

// NewAppModule creates a new AppModule object
func NewAppModule(server server.Server) AppModule {
	return AppModule{
		server: server,
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(router *sdkservice.MsgRouter) {
	configurator := router.GetConfigurator()
	// register dvs-msg handler server
	types.RegisterDVSRequestServer(configurator, am.server)

	// register dvs-msg result handler
	configurator.RegisterResultMsgExtractor(
		&types.RequestPriceFeedOut{}, resulthandlers.NewProcessRequestPriceFeedResultHandler(),
	)
}

func (am AppModule) RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}
