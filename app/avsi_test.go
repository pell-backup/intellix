package app

import (
	"cosmossdk.io/log"
	cosmossdk_io_math "cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/0xPellNetwork/pelldvs/aggregator"
	avsi "github.com/0xPellNetwork/pelldvs/application"
	"github.com/0xPellNetwork/pelldvs/avsi/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	"github.com/stretchr/testify/require"
	dvsservermanager "intellix/pkg/dvs_msg_handler"
	dvstypes "intellix/x/price/dvs/types"
	"math/big"
	"os"
	"testing"
	"time"
)

const (
	SimAppChainID = "intellix-simapp"
)

func interBlockCacheOpt() func(*baseapp.BaseApp) {
	return baseapp.SetInterBlockCache(store.NewCommitKVStoreCacheManager())
}

func mockDvsRequestData() ([]byte, error) {
	data := &dvstypes.ProcessRequestPriceFeedIn{
		TaskIndex:                 1,
		RequestId:                 []byte{},
		FeeToken:                  "01",
		Payment:                   cosmossdk_io_math.NewInt(1),
		RequestData:               nil,
		CallbackAddress:           "0x001",
		CallbackFunctionId:        []byte{},
		TaskCreatedBlock:          1,
		QuorumNumbers:             []byte{},
		QuorumThresholdPercentage: 10,
		PriceFeed: &dvstypes.PriceFeedParam{
			BaseSymbol:  "btc",
			QuoteSymbol: "usdt",
		},
	}

	return dvsservermanager.EncodeMsgs(data)
}

func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

func newApp(t *testing.T) (*App, sdk.Context) {
	simcli.FlagSeedValue = time.Now().Unix()
	simcli.FlagVerboseValue = true
	simcli.FlagCommitValue = true
	simcli.FlagEnabledValue = true

	config := simcli.NewConfigFromFlags()
	config.ChainID = SimAppChainID

	db, dir, logger, skip, err := simtestutil.SetupSimulation(config, "leveldb-app-sim", "Simulation", simcli.FlagVerboseValue, simcli.FlagEnabledValue)
	if skip {
		t.Skip("skipping application simulation")
	}
	require.NoError(t, err, "simulation setup failed")

	defer func() {
		require.NoError(t, db.Close())
		require.NoError(t, os.RemoveAll(dir))
	}()

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	bApp, err := New(logger, db, nil, true, appOptions, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
	require.NoError(t, err)
	require.Equal(t, Name, bApp.Name())

	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	storeKey := storetypes.NewKVStoreKey("dvs")
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())
	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	return bApp, ctx
}

func TestProcessRequest(t *testing.T) {
	a, ctx := newApp(t)
	data, err := mockDvsRequestData()
	if err != nil {
		t.Fatalf("error in mockDvsRequestData: %s", err.Error())
	}

	resp, err := a.ProcessRequest(ctx, &avsi.RequestProcessRequest{
		Request: types.DVSRequest{
			Data:    data,
			Height:  1,
			ChainID: big.NewInt(1),
		},
	})
	require.NoError(t, err)

	require.NotNil(t, resp.ResponseDigest)

}

func TestPostProcessRequest(t *testing.T) {
	a, ctx := newApp(t)

	reqData, err := mockDvsRequestData()
	if err != nil {
		t.Fatalf("error in mockDvsRequestData: %s", err.Error())
	}

	processResp, err := a.ProcessRequest(ctx, &avsi.RequestProcessRequest{
		Request: types.DVSRequest{
			Data:    reqData,
			Height:  1,
			ChainID: big.NewInt(1),
		},
	})
	require.NoError(t, err)

	_, err = a.PostRequest(ctx, &avsi.RequestPostRequest{
		Request: types.DVSRequest{
			Data:    reqData,
			Height:  1,
			ChainID: big.NewInt(1),
		},
		Response: aggregator.ValidatedResponse{
			Data: processResp.Reponse,
		},
	})
}
