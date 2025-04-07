package types

import (
	"github.com/0xPellNetwork/pelldvs/avsi/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	msgservice.RegisterMsgServiceDesc(registry, &DVSRequest_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &types.DVSResponse{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &RequestPriceFeedIn{})
}
