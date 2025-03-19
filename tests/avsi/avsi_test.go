package avsi

// CI: remove test
/*
const (
	SimAppChainID = "intellix-simapp"
)


func mockDvsRequestData() ([]byte, error) {
	data := &dvstypes.ProcessRequestPriceFeedIn{
		TaskMetadata: &dvstypes.TaskRequest{
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

	resp, err := a.ProcessDVSRequest(ctx, &avsi.RequestProcessDVSRequest{
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

	processResp, err := a.ProcessDVSRequest(ctx, &avsi.RequestProcessDVSRequest{
		Request: types.DVSRequest{
			Data:    reqData,
			Height:  1,
			ChainID: big.NewInt(1),
		},
	})
	require.NoError(t, err)

	_, err = a.ProcessDVSResponse(ctx, &avsi.RequestProcessDVSResponse{
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

func fauxMerkleModeOpt(bapp *baseapp.BaseApp) {
	bapp.SetFauxMerkleMode()
}

func newApp(t *testing.T) (*app.App, context.Context) {
	simcli.FlagSeedValue = time.Now().Unix()
	simcli.FlagVerboseValue = true
	simcli.FlagCommitValue = true
	simcli.FlagEnabledValue = true

	config := simcli.NewConfigFromFlags()
	config.ChainID = SimAppChainID
	config.DBBackend = "memdb"

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
	appOptions[flags.FlagHome] = "intellix"
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue

	bApp, err := app.New(logger, db, nil, true, appOptions, fauxMerkleModeOpt, baseapp.SetChainID(SimAppChainID))
	require.NoError(t, err, "app New failed")

	return bApp, context.Background()
}

func mockDvsRequestData() ([]byte, error) {
	data := &processordvstypes.RequestScriptIn{
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
		ScriptId:                  1,
		ScriptParam:               []byte("123"),
	}

	return dvsservermanager.EncodeMsgs(data)
}

// #cgo LDFLAGS: -L${SRCDIR}/../../lib -lruntime
func TestProcessorRequest(t *testing.T) {
	bApp, _ := newApp(t)
	dvsservermanager.InitDvsMsgHelper(bApp.AppCodec())

	a := &pellapp.App{}
	a.SetAppCodec(bApp.AppCodec())
	a.SetInterfaceRegistry(bApp.InterfaceRegistry())
	a.RegisterInterfaceByParam(bApp.InterfaceRegistry())
	m := processordvs.NewAppModule(a.ProcessorDvsServer)
	m.RegisterServices()
	processordvstypes.RegisterInterfaces(a.InterfaceRegistry())

	_, err := mockDvsRequestData()
	if err != nil {
		t.Fatalf("error in mockDvsRequestData: %s", err.Error())
	}

	//resp, err := a.ProcessDVSRequest(ctx, &avsitypes.RequestProcessDVSRequest{
	//	Request: &avsitypes.DVSRequest{
	//		Data:                      data,
	//		Height:                    1,
	//		ChainId:                   1,
	//		GroupNumbers:              []uint32{},
	//		GroupThresholdPercentages: []uint32{},
	//	},
	//})
	//require.NoError(t, err)
	//require.NotNil(t, resp.ResponseDigest)

}
*/
