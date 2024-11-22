package app

import (
	clienthelpers "cosmossdk.io/client/v2/helpers"
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
	*dvsconfig.Config
	RootDir        string  `mapstructure:"root_dir"`
	OperatorAddr   string  `mapstructure:"operator_addr"`
	WaitBlockCount int64   `mapstructure:"wait_block_count"`
	GasPrices      string  `mapstructure:"gas_prices"`
	GasAdjustment  float64 `mapstructure:"gas_adjustment"`
}

func (p PellAppConfig) Validate() error {
	if p.RootDir == "" {
		home, err := clienthelpers.GetNodeHomeDirectory("." + Name)
		if err != nil {
			return err
		}
		p.RootDir = home
	}
	if p.OperatorAddr == "" {
		return fmt.Errorf("no operator address provided")
	}
	if p.Config == nil || p.Config.ValidateBasic() != nil {
		return fmt.Errorf("invalid pell config")
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

	if config.RootDir == "" {
		config.RootDir = DefaultNodeHome
	}

	config.Config.RootDir = config.RootDir
	config.Config.SetRoot(config.RootDir)
	dvsconfig.EnsureRoot(config.RootDir)

	// dvs client
	var err error
	app.dvsNode, err = pelldvs.NewNode(app.logger, app, config.Config)
	if err != nil {
		panic(err)
	}

	//dvs server manager
	app.DvsServer = dvsserver.NewServer(app.logger, clientCtx, config.OperatorAddr, config.WaitBlockCount, config.GasPrices, config.GasAdjustment)
	dvsservermanager.InitDvsMsgHelper(app.appCodec)
	app.PostProcessRequestServer = dvsservermanager.GetPostProcessRequestHandler()
	app.ProcessRequestServer = dvsservermanager.GetProcessRequestHandler()
	dvs.NewAppModule(app.DvsServer).RegisterServices()
	dvstypes.RegisterInterfaces(app.interfaceRegistry)

	return app
}
