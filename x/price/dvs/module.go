package dvs

import (
	"encoding/json"
	"fmt"
	modulev1 "intellix/api/intellix/intellix/module"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	"intellix/x/price/dvs/keeper"
	dvstypes "intellix/x/price/dvs/types"
	"intellix/x/price/types"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var (
	_ appmodule.AppModule   = (*AppModule)(nil)
	_ module.AppModuleBasic = (*AppModule)(nil)
)

// AppModuleBasic defines the basic application module used by the dvs module.
type AppModuleBasic struct{}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {
}

// Name returns the dvs module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the dvs module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {}

// RegisterInterfaces registers the module's interface types
func (AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	dvstypes.RegisterInterfaces(reg)
}

// DefaultGenesis returns default genesis state as raw bytes for the dvs module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

// ValidateGenesis performs genesis state validation for the dvs module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, config client.TxEncodingConfig, bz json.RawMessage) error {
	var genState types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &genState); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}
	return genState.Validate()
}

// AppModule implements an application module for the dvs module.
type AppModule struct {
	AppModuleBasic
	keeper                   keeper.Keeper
	ProcessRequestServer     grpc1.Server
	PostProcessRequestServer grpc1.Server
}

func (am AppModule) IsOnePerModuleType() {
}

func (am AppModule) IsAppModule() {
}

func (am AppModule) RegisterGRPCGatewayRoutes(context client.Context, mux *runtime.ServeMux) {
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper, cdc codec.Codec) AppModule {
	return AppModule{
		AppModuleBasic:           AppModuleBasic{},
		keeper:                   k,
		ProcessRequestServer:     dvsservermanager.GetProcessRequestHandler(),
		PostProcessRequestServer: dvsservermanager.GetPostProcessRequestHandler(),
	}
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	dvsProcessRequestServer := keeper.NewDvsProcessRequestServer(am.keeper)
	dvsPostProcessRequestServer := keeper.NewDvsPostProcessRequestServer(am.keeper)

	// register cosmos-sdk handler server
	dvstypes.RegisterDvsProcessRequestServer(cfg.MsgServer(), dvsProcessRequestServer)
	dvstypes.RegisterDvsPostProcessRequestServer(cfg.MsgServer(), dvsPostProcessRequestServer)

	// register dvs-msg handler server
	dvstypes.RegisterDvsProcessRequestServer(am.ProcessRequestServer, dvsProcessRequestServer)
	dvstypes.RegisterDvsPostProcessRequestServer(am.PostProcessRequestServer, dvsPostProcessRequestServer)
}

func init() {
	appmodule.Register(
		&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	StoreService store.KVStoreService
	Cdc          codec.Codec
	Logger       log.Logger

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper

	clientCtx client.Context
}

type ModuleOutputs struct {
	depinject.Out

	PriceKeeper keeper.Keeper
	Module      appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// TODO: configurable operator and so on
	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.Logger,
		in.clientCtx,
		"",
		10,
		"",
		0,
	)
	m := NewAppModule(k, in.Cdc)

	return ModuleOutputs{PriceKeeper: k, Module: m}
}
