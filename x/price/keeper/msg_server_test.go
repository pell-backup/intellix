package keeper_test

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "intellix/testutil/keeper"
	"intellix/x/price/keeper"
	"intellix/x/price/types"
)

func setupMsgServer(t testing.TB) (keeper.Keeper, types.MsgServer, context.Context) {
	k, ctx := keepertest.PriceKeeper(t)
	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
	require.NotEmpty(t, k)
}

func TestMsgServer_VoteRequestPriceFeed(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)
	_, err := ms.VoteRequestPriceFeed(ctx, &types.MsgVoteRequestPriceFeed{
		TaskIndex:   1,
		OperatorId:  "operator",
		RequestId:   []byte("request_id"),
		BaseSymbol:  "usdt",
		QuoteSymbol: "btc",
		Prices: []*types.VoteRequestPriceFeed{{
			Source: "source",
			Price:  math.LegacyNewDec(100000000),
		}},
		Timestamp:   time.Now().Unix(),
		BlockHeight: 1000,
	})
	require.NoError(t, err)
}

func TestMsgServer_VoteFinalizedRequestPrice(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)
	_, err := ms.VoteFinalizedRequestPrice(ctx, &types.MsgVoteFinalizedRequestPrice{
		TaskIndex:                 1,
		RequestId:                 []byte("request_id"),
		FeeToken:                  "",
		Payment:                   math.Int{},
		RequestData:               nil,
		CallbackAddress:           "",
		CallbackFunctionId:        nil,
		TaskCreatedBlock:          0,
		QuorumNumbers:             nil,
		QuorumThresholdPercentage: 0,
		ReferenceTaskIndex:        0,
		Price:                     math.Int{},
		Sender:                    "",
	})
	require.NoError(t, err)
}
