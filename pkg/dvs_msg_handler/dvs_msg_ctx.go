package dvsservermanager

import (
	avsiTypes "github.com/0xPellNetwork/pelldvs/avsi/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ctxDvsRequestKey = "CTX_DVS_REQUEST"
)

func CtxWithDvsRequestData(ctx sdk.Context, request *avsiTypes.DVSRequest) sdk.Context {
	return ctx.WithValue(ctxDvsRequestKey, request.Data)
}

func CtxGetDvsRequestData(ctx sdk.Context) ([]byte, bool) {
	value := ctx.Value(ctxDvsRequestKey)
	val, ok := value.([]byte)
	return val, ok
}
