package cmd

import (
	"fmt"
	taskdispatcher "intellix/dispatcher"
	pkglogger "intellix/pkg/logger"
	"intellix/pkg/pelldvs"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

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
