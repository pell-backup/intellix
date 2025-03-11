package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	keepertest "intellix/testutil/keeper"
	"intellix/x/processor/types"
)

func TestListProcessor(t *testing.T) {
	k, ctx := keepertest.ProcessorKeeper(t)
	_, err := k.ListProcessor(ctx, nil)
	require.Error(t, err)
	list, err := k.ListProcessor(ctx, &types.QueryListProcessorRequest{})
	require.NoError(t, err)
	require.Len(t, list.Processor, 0)
	id, err := k.AppendProcessor(ctx, types.Processor{
		Creator:       "creator",
		ProcessorType: types.ProcessorType_WASM,
		Config:        []byte("config"),
		WasmCode:      []byte("wasm code"),
	})
	require.NoError(t, err)
	list, err = k.ListProcessor(ctx, &types.QueryListProcessorRequest{Pagination: &query.PageRequest{}})
	require.NoError(t, err)
	require.Len(t, list.Processor, 1)
	require.Equal(t, id, list.Processor[0].Id)
	require.Equal(t, "creator", list.Processor[0].Creator)
	require.Equal(t, types.ProcessorType_WASM, list.Processor[0].ProcessorType)
	require.Equal(t, []byte("config"), list.Processor[0].Config)
	require.Equal(t, []byte("wasm code"), list.Processor[0].WasmCode)
	require.Len(t, list.Pagination.NextKey, 0)
}
