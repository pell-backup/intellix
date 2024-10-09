package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "intellix/testutil/keeper"
	"intellix/x/price/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.PriceKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
