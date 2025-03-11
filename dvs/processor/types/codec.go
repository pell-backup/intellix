package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"

	dvstypes "intellix/sdk/pelldvs/types"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	msgservice.RegisterMsgServiceDesc(registry, &_DVSRequest_serviceDesc)
	msgservice.RegisterMsgServiceDesc(registry, &_DVSResponse_serviceDesc)

	registry.RegisterImplementations((*sdk.Msg)(nil), &dvstypes.RequestPostRequestValidatedData{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &RequestScriptIn{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &ResponseScriptOut{})
}
