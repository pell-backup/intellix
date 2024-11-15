package app

import (
	"context"
	"cosmossdk.io/log"
	cosmossdk_io_math "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	"github.com/stretchr/testify/require"
	"intellix/x/price/dvs/types"
	"os"
	"testing"
	"time"
)

func newBaseApp(t *testing.T) (*App, context.Context) {
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

	//stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	//storeKey := storetypes.NewKVStoreKey("dvs")
	//stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	//require.NoError(t, stateStore.LoadLatestVersion())
	//ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger())

	return bApp, context.Background()
}

func TestNewPellApp(t *testing.T) {
	app, _ := newBaseApp(t)
	var testMsg = &types.TaskRequest{
		TaskIndex:                 1,
		RequestId:                 []byte("1"),
		FeeToken:                  "1",
		Payment:                   cosmossdk_io_math.Int{},
		RequestData:               []byte("btc-usdt"),
		CallbackAddress:           "0x10101",
		CallbackFunctionId:        nil,
		TaskCreatedBlock:          0,
		QuorumNumbers:             nil,
		QuorumThresholdPercentage: 0,
	}

	require.NotPanics(t,
		func() {
			papp := NewPellApp(log.NewLogger(os.Stdout), &PellAppConfig{})

			ret1, err := papp.appCodec.Marshal(testMsg)
			require.NoError(t, err)

			ret2, err := app.appCodec.Marshal(testMsg)
			require.NoError(t, err)

			require.Equal(t, ret1, ret2)
		},
	)
}
