package app

const (
	SimAppChainID = "intellix-simapp"
)

// CI: remove test
/*

func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

func mockDvsRequestData() ([]byte, error) {
	data := &dvstypes.ProcessRequestPriceFeedIn{
		Task: &dvstypes.TaskRequest{
			TaskIndex:                 1,
			RequestId:                 []byte("1234"),
			FeeToken:                  "",
			Payment:                   cosmossdk_io_math.NewInt(1),
			RequestData:               []byte(""),
			CallbackAddress:           "",
			CallbackFunctionId:        []byte("1"),
			TaskCreatedBlock:          2,
			QuorumNumbers:             []byte("1"),
			QuorumThresholdPercentage: 10,
		},
		PriceFeed: &dvstypes.PriceFeedParam{
			BaseSymbol:  "btc",
			QuoteSymbol: "usdt",
		},
	}

	return dvsservermanager.EncodeMsgs(data)
}

func newApp(t *testing.T) (*App, context.Context) {
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

	_ = NewPellApp(bApp.logger, &PellAppConfig{})
	return bApp, context.Background()
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
*/
