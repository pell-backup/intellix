package cmd

import (
	"fmt"

	"cosmossdk.io/errors"
	pelldvscfg "github.com/0xPellNetwork/pelldvs/config"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	dispatcher "intellix/dispatcher"
	sdklogger "intellix/sdk/logger"
)

func startDispatcherCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-dispatcher",
		Short: "Start the dispatcher service",
		Long: "Start the dispatcher service, Example:\n" +
			"intellixd start-dispatcher",
		RunE: func(cmd *cobra.Command, args []string) error {
			serverCtx := server.GetServerContextFromCmd(cmd)

			logger := sdklogger.NewDVSLogAdapter(serverCtx.Logger)

			home := getConfigHome()

			// load dispatcher config
			var conf = &dispatcher.Config{}
			viper.SetConfigFile(home + "/config/dispatcher.config.json")
			if err := viper.ReadInConfig(); err != nil {
				return err
			}
			if err := viper.Unmarshal(conf); err != nil {
				return err
			}
			if err := conf.Validate(); err != nil {
				return err
			}

			// load pellDVS config
			pellDVSConf := pelldvscfg.DefaultConfig()
			vp := viper.New()
			vp.SetConfigFile(home + "/config/config.toml")

			if err := vp.ReadInConfig(); err != nil {
				return errors.Wrap(err, "failed to read in pelldvs config")
			}
			if err := vp.Unmarshal(pellDVSConf); err != nil {
				return errors.Wrap(err, "failed to unmarshal pelldvs configuration")
			}

			pellDVSConf.SetRoot(home)
			logger.Info("PellDVS configuration",
				"config", fmt.Sprintf("%+v", pellDVSConf),
				"pell", fmt.Sprintf("%+v", pellDVSConf.Pell),
			)

			// create TaskDispatcher
			td, err := dispatcher.NewTaskDispatcher(logger, pellDVSConf, nil)
			if err != nil {
				return fmt.Errorf("failed to create TaskDispatcher: %w", err)
			}

			if err = td.Start(); err != nil {
				return fmt.Errorf("failed to start TaskDispatcher: %w", err)
			}

			<-td.Quit()

			return nil
		},
	}
	return cmd
}
