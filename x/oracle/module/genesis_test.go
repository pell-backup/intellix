package oracle_test

import (
	"testing"

	keepertest "intellix/testutil/keeper"
	"intellix/testutil/nullify"
	oracle "intellix/x/oracle/module"
	"intellix/x/oracle/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.OracleKeeper(t)
	oracle.InitGenesis(ctx, k, genesisState)
	got := oracle.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
