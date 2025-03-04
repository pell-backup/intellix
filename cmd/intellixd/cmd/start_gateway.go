package cmd

import (
	"context"
	"fmt"

	taskgateway "intellix/gateway"
	pkglogger "intellix/sdk/logger"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

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

			dvsLogger.Info("Starting TaskGateway service", "config", conf)
			taskGateway, err := taskgateway.NewTaskGateway(dvsLogger, context.Background(), conf)
			if err != nil {
				return fmt.Errorf("failed to create TaskGateway: %w", err)
			}

			dvsLogger.Info("prepare to start TaskGateway service")
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
