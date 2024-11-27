package cmd

import (
	"context"
	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	"errors"
	"fmt"
	dvsconfig "github.com/0xPellNetwork/pelldvs/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
	"intellix/pkg/logger"
	pkglogger "intellix/pkg/logger"
	"intellix/pkg/pelldvs"
	"intellix/pkg/taskdispatcher"
	"intellix/pkg/taskgateway"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"intellix/app"
)

var (
	configFile string
	Name       = "intellix"
)

func initRootCmd(
	rootCmd *cobra.Command,
	txConfig client.TxConfig,
	basicManager module.BasicManager,
) {
	rootCmd.AddCommand(
		genutilcli.InitCmd(basicManager, app.DefaultNodeHome),
		NewInPlaceTestnetCmd(addModuleInitFlags),
		debug.Cmd(),
		confixcmd.ConfigCommand(),
		pruning.Cmd(newApp, app.DefaultNodeHome),
		snapshot.Cmd(newApp),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

	// add keybase, auxiliary RPC, query, genesis, and tx child commands
	rootCmd.AddCommand(
		server.StatusCommand(),
		genesisCommand(txConfig, basicManager),
		queryCommand(),
		txCommand(),
		keys.Commands(),
		taskDispatcherCommand(),
		taskGatewayCommand(),
		pellAppCommand(),
	)
}

func addModuleInitFlags(startCmd *cobra.Command) {
	crisis.AddModuleInitFlags(startCmd)
}

// genesisCommand builds genesis-related `intellixd genesis` command. Users may provide application specific commands as a parameter
func genesisCommand(txConfig client.TxConfig, basicManager module.BasicManager, cmds ...*cobra.Command) *cobra.Command {
	cmd := genutilcli.Commands(txConfig, basicManager, app.DefaultNodeHome)

	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		rpc.ValidatorCommand(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
	)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

// newApp creates the application
func newApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)

	a, err := app.New(
		logger, db, traceStore, true,
		appOpts,
		baseappOptions...,
	)
	if err != nil {
		panic(err)
	}
	return a
}

// appExport creates a new app (optionally at a given height) and exports state.
func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	var (
		bApp *app.App
		err  error
	)

	// this check is necessary as we use the flag in x/upgrade.
	// we can exit more gracefully by checking the flag here.
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, errors.New("application home not set")
	}

	viperAppOpts, ok := appOpts.(*viper.Viper)
	if !ok {
		return servertypes.ExportedApp{}, errors.New("appOpts is not viper.Viper")
	}

	// overwrite the FlagInvCheckPeriod
	viperAppOpts.Set(server.FlagInvCheckPeriod, 1)
	appOpts = viperAppOpts

	if height != -1 {
		bApp, err = app.New(logger, db, traceStore, false, appOpts)
		if err != nil {
			return servertypes.ExportedApp{}, err
		}

		if err := bApp.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	} else {
		bApp, err = app.New(logger, db, traceStore, true, appOpts)
		if err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return bApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}

// TODO: put start logic into "start" command with flag
// taskDispatcherCommand builds task-dispatcher command
func taskDispatcherCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-task-dispatcher",
		Short: "Start the TaskDispatcher service",
		Long: "Start the TaskDispatcher service, Example:\n" +
			"intellixd start-task-dispatcher --config=config.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			home := getConfigHome()
			if configFile == "" {
				configFile = home + "/config/dispatcher.config.json"
			}

			viper.SetConfigFile(configFile)
			if err := viper.ReadInConfig(); err != nil {
				return err
			}
			if err := viper.Unmarshal(config); err != nil {
				return err
			}

			var conf = &taskdispatcher.Config{}
			err := viper.Unmarshal(conf)
			if err != nil {
				return err
			}
			if err := conf.Validate(); err != nil {
				return err
			}

			// start task dispatcher
			dvsLogger := pkglogger.NewDVSLogAdapter(serverCtx.Logger)

			// new pell-dvs client
			pellDVSClient, err := pelldvs.NewClient(dvsLogger.With("module", "client"), conf.DvsAddress)
			if err != nil {
				return fmt.Errorf("failed to create PellDVS client: %w", err)
			}

			td, err := taskdispatcher.NewTaskDispatcher(dvsLogger.With("module", "task-dispacther"), pellDVSClient, conf.Chains)
			if err != nil {
				return fmt.Errorf("failed to create TaskDispatcher: %w", err)
			}

			err = td.Start()
			if err != nil {
				return fmt.Errorf("failed to start TaskDispatcher: %w", err)
			}

			// wait for quit signal
			<-td.Quit()

			return nil
		},
	}

	// add config flag
	cmd.Flags().StringVar(&configFile, "config", "", "config file")

	return cmd
}

func taskGatewayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-task-gateway",
		Short: "Start the TaskGateway service",
		Long: "Start the TaskGateway service, Example:\n" +
			"intellixd start-task-gateway --config=config.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			home := getConfigHome()
			if configFile == "" {
				configFile = home + "/config/gateway.config.json"
			}
			viper.SetConfigFile(configFile)
			if err := viper.ReadInConfig(); err != nil {
				return err
			}

			conf := &taskgateway.TaskGatewayCfg{}
			err := viper.Unmarshal(conf)
			if err != nil {
				return err
			}
			err = conf.Validate()
			if err != nil {
				return err
			}

			dvsLogger := pkglogger.NewDVSLogAdapter(serverCtx.Logger)
			taskGateway, err := taskgateway.NewTaskGateway(dvsLogger, context.Background(), conf)
			if err != nil {
				return fmt.Errorf("failed to create TaskGateway: %w", err)
			}

			err = taskGateway.Start()
			if err != nil {
				return fmt.Errorf("failed to start TaskGateway: %w", err)
			}
			<-taskGateway.Quit()

			return nil
		},
	}
	cmd.Flags().StringVar(&configFile, "config", "", "config file")
	return cmd
}

func pellAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-operator",
		Short: "Start the PellApp Operator service",
		Long: "Start the PellApp Operator service, Example:\n" +
			"intellixd start-operator --config=config.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)
			home := getConfigHome()

			if configFile == "" {
				configFile = home + "/config/operator.config.json"
			}

			vp := viper.New()

			vp.SetConfigFile(configFile)
			if err := vp.ReadInConfig(); err != nil {
				panic(err)
			}

			var pellAppConfig = &app.PellAppConfig{}
			err := vp.Unmarshal(pellAppConfig)
			if err != nil {
				panic(err)
			}
			if pellAppConfig.RootDir == "" {
				pellAppConfig.RootDir = home
			}

			// read pellConfig
			if pellAppConfig.DvsConfig == nil || pellAppConfig.DvsConfig.Pell == nil {
				pellAppConfig.DvsConfig = dvsconfig.DefaultConfig()
				if _, err := os.Stat(home + "/config/config.toml"); err == nil {
					vp.SetConfigFile(home + "/config/config.toml")
					if err := vp.ReadInConfig(); err != nil {
						return err
					}
					if err := vp.Unmarshal(pellAppConfig.DvsConfig); err != nil {
						return err
					}
				}
			}

			err = pellAppConfig.Validate()
			if err != nil {
				panic(err)
			}

			dApp, err := newDefaultApp(getConfigHome(), serverCtx.Logger)
			if err != nil {
				panic(err)
			}

			a := app.NewPellApp(dApp.InterfaceRegistry(),
				logger.NewDVSLogAdapter(serverCtx.Logger), pellAppConfig,
			)
			return a.Start()
		},
	}
	cmd.Flags().StringVar(&configFile, "config", "", "config file")
	return cmd
}

func getConfigHome() string {
	home := os.Getenv("PELLDVS_HOME")
	if home == "" {
		home, _ = clienthelpers.GetNodeHomeDirectory("." + Name)
	}
	return home
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	return dbm.NewDB("application", dbm.GoLevelDBBackend, dataDir)
}

func newDefaultApp(home string, log log.Logger) (*app.App, error) {
	db, err := openDB(home)
	if err != nil {
		return nil, err
	}
	newVp := viper.New()
	newVp.Set("pruning", "default")
	newVp.Set("home", getConfigHome())
	//newVp.Set("chain-id", pellAppConfig.CosmosChainId)
	newVp.Set("chain-id", "intellix")
	baseApp := newApp(log, db, nil, newVp)
	return baseApp.(*app.App), nil
}
