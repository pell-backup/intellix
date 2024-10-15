package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&ProcessPriceFeedMsg{},
		&ResponsePostRequestPriceFeed{},
		&PostPriceFeedMsg{},
		&ResponsePostRequestPriceFeed{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Dvs_serviceDesc)
}
