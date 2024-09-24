package intellix_test

import (
	"testing"

	keepertest "intellix/testutil/keeper"
	"intellix/testutil/nullify"
	intellix "intellix/x/intellix/module"
	"intellix/x/intellix/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.IntellixKeeper(t)
	intellix.InitGenesis(ctx, k, genesisState)
	got := intellix.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
