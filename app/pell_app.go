package app

import (
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	"intellix/x/price/dvs"
	dvsserver "intellix/x/price/dvs/server"
	dvstypes "intellix/x/price/dvs/types"
)

type PellApp struct {
	logger log.Logger

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	DvsServer                dvsserver.Server
	ProcessRequestServer     grpc1.Server
	PostProcessRequestServer grpc1.Server
}

type PellAppConfig struct {
	OperatorAddr   string  `mapstructure:"operator_addr"`
	WaitBlockCount int64   `mapstructure:"wait_block_count"`
	GasPrices      string  `mapstructure:"gas_prices"`
	GasAdjustment  float64 `mapstructure:"gas_adjustment"`
}

func NewPellApp(
	logger log.Logger,
	config *PellAppConfig,
) *PellApp {
	var app = &PellApp{}

	appConfig = depinject.Configs(
		depinject.Configs(
			appConfig,
		),
		depinject.Supply(
			map[string]interface{}{}, // supply app options
			logger,                   // supply logger
		),
	)

	if err := depinject.Inject(appConfig,
		&app.logger,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
	); err != nil {
		panic(err)
	}

	clientCtx := NewClientContext(app.appCodec, app.interfaceRegistry, tx.ConfigOptions{}, app.legacyAmino)

	//dvs server manager
	app.DvsServer = dvsserver.NewServer(app.logger, clientCtx, config.OperatorAddr, config.WaitBlockCount, config.GasPrices, config.GasAdjustment)
	dvsservermanager.InitDvsMsgHelper(app.appCodec)
	app.PostProcessRequestServer = dvsservermanager.GetPostProcessRequestHandler()
	app.ProcessRequestServer = dvsservermanager.GetProcessRequestHandler()
	dvs.NewAppModule(app.DvsServer).RegisterServices()
	dvstypes.RegisterInterfaces(app.interfaceRegistry)

	return app
}
