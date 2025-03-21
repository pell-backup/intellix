package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &VRFMsgRequest_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &VRFTaskRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &VRFTaskResponse{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &DVSResultResponse{})
}
