package price_test

import (
	"testing"

	keepertest "intellix/testutil/keeper"
	"intellix/testutil/nullify"
	price "intellix/x/price/module"
	"intellix/x/price/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.PriceKeeper(t)
	price.InitGenesis(ctx, k, genesisState)
	got := price.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
