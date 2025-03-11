package processor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "intellix/testutil/keeper"
	"intellix/testutil/nullify"
	processor "intellix/x/processor/module"
	"intellix/x/processor/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ProcessorKeeper(t)
	processor.InitGenesis(ctx, k, genesisState)
	got := processor.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
