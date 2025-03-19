package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &VRFMsgServer_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &GenerateRandomNumberRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &GenerateRandomNumberResponse{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &DVSResultResponse{})
}
