package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil), &ProcessPriceFeedMsg{})
	msgservice.RegisterMsgServiceDesc(registry, &_DvsProcessRequest_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &RequestPostRequestPriceFeedValidatedData{})
	msgservice.RegisterMsgServiceDesc(registry, &_DvsPostProcessRequest_serviceDesc)
}
