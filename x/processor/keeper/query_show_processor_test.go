package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "intellix/testutil/keeper"
	"intellix/x/processor/types"
)

func TestShowProcessor(t *testing.T) {
	k, ctx := keepertest.ProcessorKeeper(t)
	_, err := k.ShowProcessor(ctx, nil)
	require.Error(t, err)
	_, err = k.ShowProcessor(ctx, &types.QueryShowProcessorRequest{Id: 0})
	require.Error(t, err)
	id, err := k.AppendProcessor(ctx, types.Processor{
		Creator:       "creator",
		ProcessorType: types.ProcessorType_PROCESSOR_TYPE_WASM,
		Config:        []byte("config"),
		WasmCode:      []byte("wasm code"),
	})
	require.NoError(t, err)
	processor, err := k.ShowProcessor(ctx, &types.QueryShowProcessorRequest{Id: id})
	require.NoError(t, err)
	require.Equal(t, id, processor.Processor.Id)
	require.Equal(t, "creator", processor.Processor.Creator)
	require.Equal(t, types.ProcessorType_PROCESSOR_TYPE_WASM, processor.Processor.ProcessorType)
	require.Equal(t, []byte("config"), processor.Processor.Config)
	require.Equal(t, []byte("wasm code"), processor.Processor.WasmCode)
}
