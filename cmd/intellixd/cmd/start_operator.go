package cmd

import (
	"fmt"
	taskdispatcher "intellix/dispatcher"
	"intellix/sdk/logger"
	pkglogger "intellix/sdk/logger"
	"os"

	dvsconfig "github.com/0xPellNetwork/pelldvs/config"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	app "intellix/pellapp"
)

var configDispatcherFile string

func pellAppCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-operator",
		Short: "Start the PellApp Operator service",
		Long: "Start the PellApp Operator service, Example:\n" +
			"intellixd start-operator--config=config.yml --config-dispatcher=config.yaml",
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

			var pellAppConfig = &app.AppConfig{}
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

			a := app.NewApp(dApp.InterfaceRegistry(),
				logger.NewDVSLogAdapter(serverCtx.Logger), pellAppConfig,
			)

			config := serverCtx.Config
			if configDispatcherFile == "" {
				configDispatcherFile = home + "/config/dispatcher.config.json"
			}

			viper.SetConfigFile(configDispatcherFile)
			if err := viper.ReadInConfig(); err != nil {
				return err
			}
			if err := viper.Unmarshal(config); err != nil {
				return err
			}

			var conf = &taskdispatcher.Config{}
			err = viper.Unmarshal(conf)
			if err != nil {
				return err
			}
			if err := conf.Validate(); err != nil {
				return err
			}

			// start task dispatcher
			dvsLogger := pkglogger.NewDVSLogAdapter(serverCtx.Logger)

			td, err := taskdispatcher.NewTaskDispatcher(dvsLogger.With("module", "task-dispacther"), a.DVSClient, conf.Chains)
			if err != nil {
				return fmt.Errorf("failed to create TaskDispatcher: %w", err)
			}

			err = td.Start()
			if err != nil {
				return fmt.Errorf("failed to start TaskDispatcher: %w", err)
			}

			//start Operator
			if err = a.Start(); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&configFile, "config", "", "config file")
	cmd.Flags().StringVar(&configDispatcherFile, "config-dispatcher", "", "dispatcher config file")

	return cmd
}
