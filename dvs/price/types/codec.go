package types

import (
	dvstypes "intellix/pkg/pelldvs/types"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	msgservice.RegisterMsgServiceDesc(registry, &_DVSRequest_serviceDesc)
	msgservice.RegisterMsgServiceDesc(registry, &_DVSResponse_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &dvstypes.RequestPostRequestValidatedData{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &RequestPriceFeedIn{})

}
