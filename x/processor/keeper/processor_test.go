package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "intellix/testutil/keeper"
	"intellix/x/processor/types"
)

func TestGetProcessorCount(t *testing.T) {
	k, ctx := keepertest.ProcessorKeeper(t)
	require.Equal(t, uint64(0), k.GetProcessorCount(ctx))
}

func TestSetProcessorCount(t *testing.T) {
	k, ctx := keepertest.ProcessorKeeper(t)
	k.SetProcessorCount(ctx, 3)
	require.Equal(t, uint64(3), k.GetProcessorCount(ctx))
}

func processor() types.Processor {
	return types.Processor{
		ProcessorType: types.ProcessorType_PROCESSOR_TYPE_WASM,
		Config:        []byte("config"),
		WasmCode:      []byte("wasm code"),
		Creator:       "creator",
	}
}

func TestAppendProcessor(t *testing.T) {
	k, ctx := keepertest.ProcessorKeeper(t)
	id, err := k.AppendProcessor(ctx, processor())
	require.NoError(t, err)
	require.Equal(t, uint64(0), id)
	require.Equal(t, uint64(1), k.GetProcessorCount(ctx))
}

func TestGetProcessor(t *testing.T) {
	k, ctx := keepertest.ProcessorKeeper(t)
	id, err := k.AppendProcessor(ctx, processor())
	require.NoError(t, err)
	processor, found := k.GetProcessor(ctx, id)
	require.True(t, found)
	require.Equal(t, processor.ProcessorType, types.ProcessorType_PROCESSOR_TYPE_WASM)
	require.Equal(t, processor.Config, []byte("config"))
	require.Equal(t, processor.WasmCode, []byte("wasm code"))
	require.Equal(t, processor.Creator, "creator")
}
