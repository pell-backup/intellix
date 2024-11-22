package app

import (
	"fmt"
	dvsconfig "github.com/0xPellNetwork/pelldvs/config"
	"github.com/0xPellNetwork/pelldvs/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	"intellix/pkg/pelldvs"
	"intellix/x/price/dvs"
	dvsserver "intellix/x/price/dvs/server"
	dvstypes "intellix/x/price/dvs/types"
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

	RootDir      string `mapstructure:"root_dir"`
	GatewayAddr  string `mapstructure:"gateway_addr"`
	OperatorAddr string `mapstructure:"operator_address"`

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
	return nil
}

func (p *PellApp) InterfaceRegistry() codectypes.InterfaceRegistry {
	if p.interfaceRegistry == nil {
		p.interfaceRegistry = codectypes.NewInterfaceRegistry()
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
	if err := p.dvsNode.Start(); err != nil {
		return err
	}
	c := make(chan interface{})
	<-c
	return nil
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
	clientCtx := NewClientContext(app.appCodec, app.interfaceRegistry, tx.ConfigOptions{})

	config.DvsConfig.RootDir = config.RootDir
	config.DvsConfig.SetRoot(config.RootDir)
	dvsconfig.EnsureRoot(config.DvsConfig.RootDir)

	// dvs client
	var err error
	app.dvsNode, err = pelldvs.NewNode(app.logger, app, config.DvsConfig)
	if err != nil {
		panic(err)
	}

	//dvs server manager
	app.DvsServer, err = dvsserver.NewServer(
		app.logger, clientCtx, config.GatewayAddr, config.OperatorAddr,
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
