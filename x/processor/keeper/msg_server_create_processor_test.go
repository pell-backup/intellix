package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "intellix/testutil/keeper"
	"intellix/x/processor/keeper"
	"intellix/x/processor/types"
)

func TestCreateProcessor(t *testing.T) {
	k, ctx := keepertest.ProcessorKeeper(t)
	msgServer := keeper.NewMsgServerImpl(k)
	msg := types.MsgCreateProcessor{
		Creator:       "creator",
		ProcessorType: types.ProcessorType_WASM,
		Config:        []byte("config"),
		WasmCode:      []byte("wasm code"),
	}
	res, err := msgServer.CreateProcessor(ctx, &msg)
	require.NoError(t, err)
	require.Equal(t, uint64(0), res.Id)
}
