package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	dvstypes "intellix/pkg/pelldvs/types"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	msgservice.RegisterMsgServiceDesc(registry, &_DvsProcessRequest_serviceDesc)
	msgservice.RegisterMsgServiceDesc(registry, &_DvsPostProcessRequest_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &ProcessPriceFeedMsg{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &dvstypes.RequestPostRequestValidatedData{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &AggregatedRequestPrice{})
}
