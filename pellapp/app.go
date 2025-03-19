package pellapp

import (
	"os"

	"github.com/0xPellNetwork/pellapp-sdk/baseapp"
	"github.com/0xPellNetwork/pellapp-sdk/pelldvs"
	"github.com/0xPellNetwork/pelldvs-libs/log"
	dvsconfig "github.com/0xPellNetwork/pelldvs/config"
	rpclocal "github.com/0xPellNetwork/pelldvs/rpc/client/local"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"

	dvs "intellix/dvs/price"
	priceserver "intellix/dvs/price/server"
	processordvs "intellix/dvs/processor"
	processorserver "intellix/dvs/processor/server"
	vrfserver "intellix/dvs/vrf/server"
	"intellix/sdk/baseapp"
	"intellix/sdk/pelldvs"
)

const (
	Name = "intellix"
)

var (
	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

type App struct {
	*baseapp.BaseApp

	appCodec          codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry

	dvsNode *pelldvs.Node

	PriceServer     priceserver.Server
	ProcessorServer processorserver.Server
	VRFServer       vrfserver.Server

	ProcessRequestServer     grpc1.Server
	PostProcessRequestServer grpc1.Server

	logger log.Logger

	DVSClient *rpclocal.Local
}

func (app *App) SetInterfaceRegistry(registry codectypes.InterfaceRegistry) codectypes.InterfaceRegistry {
	app.interfaceRegistry = registry
	return app.interfaceRegistry
}

func (app *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	if app.interfaceRegistry == nil {
		app.interfaceRegistry = codectypes.NewInterfaceRegistry()
	}
	return app.interfaceRegistry
}

func (app *App) RegisterInterface() {
	std.RegisterInterfaces(app.interfaceRegistry)
	sdktypes.RegisterInterfaces(app.interfaceRegistry)
}

func (app *App) RegisterInterfaceByParam(p codectypes.InterfaceRegistry) {
	std.RegisterInterfaces(p)
	sdktypes.RegisterInterfaces(p)
}

func (app *App) SetAppCodec(codec codec.Codec) {
	app.appCodec = codec
}

func (app *App) AppCodec() codec.Codec {
	if app.appCodec == nil {
		app.appCodec = codec.NewProtoCodec(app.interfaceRegistry)
	}
	return app.appCodec
}

func (app *App) Start() error {
	app.logger.Info("App Start")
	if err := app.dvsNode.Start(); err != nil {
		app.logger.Error("DvsNode Start Failed", "error", err.Error())
		return err
	}
	c := make(chan interface{})
	<-c
	return nil
}

func getOperatorName() string {
	name := os.Getenv("OPERATOR_KEY_NAME")
	if name == "" {
		return Name
	}
	return name
}

func NewApp(
	interfaceRegistry codectypes.InterfaceRegistry,
	logger log.Logger,
	config *AppConfig,
) *App {
	cdc := codec.NewProtoCodec(interfaceRegistry)

	var app = &App{
		// TODO: use depinject
		BaseApp:           baseapp.NewBaseApp(Name, logger, cdc),
		interfaceRegistry: interfaceRegistry,
		logger:            logger,
		appCodec:          cdc,
	}

	app.RegisterInterface()
	var err error

	config.DvsConfig.RootDir = config.RootDir
	config.DvsConfig.SetRoot(config.RootDir)
	dvsconfig.EnsureRoot(config.DvsConfig.RootDir)

	// TODO: configurable
	config.DvsConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	// price client
	logger.Info("NewNode", "config.DvsConfig.RPC.ListenAddress", config.DvsConfig.RPC.ListenAddress)
	app.dvsNode, err = pelldvs.NewNode(app.logger, app, config.DvsConfig)
	if err != nil {
		panic(err)
	}

	// TODO: use rpc/Client
	app.DVSClient = app.dvsNode.GetLocalClient()

	// cosmos network
	kr, err := keyring.New(Name, keyring.BackendTest, config.RootDir, os.Stdin, app.appCodec)
	if err != nil {
		panic(err)
	}
	clientCtx := NewClientContext(app.appCodec, app.interfaceRegistry, tx.ConfigOptions{}, getOperatorName(), config.CosmosNodeUri, config.CosmosChainId, kr)

	key, err := clientCtx.Keyring.Key(getOperatorName())
	if err != nil {
		panic(err)
	}

	app.PriceServer, err = priceserver.NewServer(
		app.logger, clientCtx, key, config.CosmosChainId,
		config.GatewayAddr, config.OperatorAddr,
		config.WaitBlockCount, config.GasPrices, config.GasAdjustment, config.PriceTickConverterConfig,
	)
	if err != nil {
		panic(err)
	}

	priceModule := dvs.NewAppModule(app.PriceServer)
	priceModule.RegisterServices(app.GetMsgRouter())
	priceModule.RegisterInterfaces(app.interfaceRegistry)

	app.ProcessorServer, err = processorserver.NewServer(
		app.logger, clientCtx, key, config.CosmosChainId,
		config.GatewayAddr, config.OperatorAddr,
		config.WaitBlockCount, config.GasPrices, config.GasAdjustment,
	)
	if err != nil {
		panic(err)
	}

	app.VRFServer, err = vrfserver.NewServer(app.logger, config.GatewayAddr)
	if err != nil {
		panic(err)
	}

	processorModule := processordvs.NewAppModule(app.ProcessorServer)
	processorModule.RegisterServices(app.GetMsgRouter())
	processorModule.RegisterInterfaces(app.interfaceRegistry)

	return app
}
