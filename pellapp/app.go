package pellapp

import (
	"os"

	"github.com/0xPellNetwork/pelldvs/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"

	baseapp "intellix/avsi"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	"intellix/pkg/pelldvs"
	"intellix/x/price/dvs"
	dvsserver "intellix/x/price/dvs/server"
	dvstypes "intellix/x/price/dvs/types"

	dvsconfig "github.com/0xPellNetwork/pelldvs/config"
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

	DvsServer                dvsserver.Server
	ProcessRequestServer     grpc1.Server
	PostProcessRequestServer grpc1.Server

	logger log.Logger
}

func (app *App) InterfaceRegistry() codectypes.InterfaceRegistry {
	if app.interfaceRegistry == nil {
		app.interfaceRegistry = codectypes.NewInterfaceRegistry()
	}
	return app.interfaceRegistry
}

func (app *App) registerInterface() {
	std.RegisterInterfaces(app.interfaceRegistry)
	sdktypes.RegisterInterfaces(app.interfaceRegistry)
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
	var app = &App{
		// TODO: use depinject
		BaseApp:           baseapp.NewBaseApp(logger),
		interfaceRegistry: interfaceRegistry,
		logger:            logger,
	}

	app.registerInterface()
	app.appCodec = app.AppCodec()
	var err error

	config.DvsConfig.RootDir = config.RootDir
	config.DvsConfig.SetRoot(config.RootDir)
	dvsconfig.EnsureRoot(config.DvsConfig.RootDir)

	// TODO: configurable
	config.DvsConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"
	// dvs client
	logger.Info("NewNode", "config.DvsConfig.RPC.ListenAddress", config.DvsConfig.RPC.ListenAddress)
	app.dvsNode, err = pelldvs.NewNode(app.logger, app, config.DvsConfig)
	if err != nil {
		panic(err)
	}

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

	//dvs server manager
	app.DvsServer, err = dvsserver.NewServer(
		app.logger, clientCtx, key, config.CosmosChainId,
		config.GatewayAddr, config.OperatorAddr,
		config.WaitBlockCount, config.GasPrices, config.GasAdjustment,
	)
	if err != nil {
		panic(err)
	}
	dvsservermanager.InitDvsMsgHelper(app.appCodec)
	app.PostProcessRequestServer = dvsservermanager.GetPostProcessRequestHandler()
	app.ProcessRequestServer = dvsservermanager.GetProcessRequestHandler()
	dvs.NewAppModule(app.DvsServer).RegisterServices()
	dvstypes.RegisterInterfaces(app.interfaceRegistry)

	return app
}
