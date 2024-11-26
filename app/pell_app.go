package app

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	"intellix/pkg/pelldvs"
	"intellix/x/price/dvs"
	dvsserver "intellix/x/price/dvs/server"
	dvstypes "intellix/x/price/dvs/types"
	"os"

	dvsconfig "github.com/0xPellNetwork/pelldvs/config"
	"github.com/0xPellNetwork/pelldvs/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
)

type PellApp struct {
	logger            log.Logger
	appCodec          codec.Codec
	interfaceRegistry codectypes.InterfaceRegistry

	dvsNode *pelldvs.Node

	DvsServer                dvsserver.Server
	ProcessRequestServer     grpc1.Server
	PostProcessRequestServer grpc1.Server
}

type PellAppConfig struct {
	DvsConfig *dvsconfig.Config `mapstructure:"-"`

	RootDir       string `mapstructure:"root_dir"`
	GatewayAddr   string `mapstructure:"gateway_addr"`
	OperatorAddr  string `mapstructure:"operator_address"`
	CosmosNodeUri string `mapstructure:"cosmos_node_uri"`
	CosmosChainId string `mapstructure:"cosmos_chain_id"`

	WaitBlockCount int64   `mapstructure:"wait_block_count"`
	GasPrices      string  `mapstructure:"gas_prices"`
	GasAdjustment  float64 `mapstructure:"gas_adjustment"`
}

func (p PellAppConfig) Validate() error {
	if p.OperatorAddr == "" {
		return fmt.Errorf("no operator address provided")
	}
	if p.DvsConfig == nil || p.DvsConfig.ValidateBasic() != nil {
		return fmt.Errorf("invalid pell config")
	}
	if p.GatewayAddr == "" {
		return fmt.Errorf("no gateway address provided")
	}
	if p.CosmosNodeUri == "" {
		return fmt.Errorf("no cosmos node uri provided")
	}
	return nil
}

func (p *PellApp) InterfaceRegistry() codectypes.InterfaceRegistry {
	if p.interfaceRegistry == nil {
		p.interfaceRegistry = codectypes.NewInterfaceRegistry()
		std.RegisterInterfaces(p.interfaceRegistry)
		sdktypes.RegisterInterfaces(p.interfaceRegistry)
	}
	return p.interfaceRegistry
}

func (p *PellApp) AppCodec() codec.Codec {
	if p.appCodec == nil {
		p.appCodec = codec.NewProtoCodec(p.interfaceRegistry)
	}
	return p.appCodec
}

func (p *PellApp) Start() error {
	p.logger.Info("PellApp Start")
	if err := p.dvsNode.Start(); err != nil {
		p.logger.Error("DvsNode Start Failed", "error", err.Error())
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

func NewPellApp(
	logger log.Logger,
	config *PellAppConfig,
) *PellApp {
	var app = &PellApp{
		logger: logger,
	}

	app.interfaceRegistry = app.InterfaceRegistry()
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
