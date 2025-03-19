package cmd

import (
	"fmt"

	dvslog "github.com/0xPellNetwork/pelldvs-libs/log"
	pelldvscfg "github.com/0xPellNetwork/pelldvs/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	taskdispatcher "intellix/dispatcher"
	taskgateway "intellix/gateway"
)

func taskGatewayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-task-gateway",
		Short: "Start the TaskGateway service",
		Long: "Start the TaskGateway service, Example:\n" +
			"intellixd start-task-gateway --config=config.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			//	serverCtx := server.GetServerContextFromCmd(cmd)
			logger := dvslog.NewNopLogger()

			home := getConfigHome()
			if configFile == "" {
				configFile = home + "/config/gateway.config.json"
			}

			conf, err := taskgateway.LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load TaskGateway configuration: %w", err)
			}
			if err = conf.Validate(); err != nil {
				return err
			}

			// load pellDVS config
			pellDVSConf := pelldvscfg.DefaultConfig()
			vp := viper.New()
			vp.SetConfigFile(home + "/config/config.toml")

			if err = vp.ReadInConfig(); err != nil {
				return errors.Wrap(err, "failed to read in pelldvs config")
			}
			if err = vp.Unmarshal(pellDVSConf); err != nil {
				return errors.Wrap(err, "failed to unmarshal pelldvs configuration")
			}

			pellDVSConf.SetRoot(home)
			logger.Info("PellDVS configuration",
				"config", fmt.Sprintf("%+v", pellDVSConf),
				"pell", fmt.Sprintf("%+v", pellDVSConf.Pell),
			)

			g, ctx := errgroup.WithContext(cmd.Context())

			// create TaskDispatcher
			tdConf := &taskdispatcher.Config{Chains: conf.Chains}
			taskDispatcher, err := taskdispatcher.NewTaskDispatcher(logger, pellDVSConf, tdConf)
			if err != nil {
				return fmt.Errorf("failed to create TaskDispatcher: %w", err)
			}

			// create TaskGateway
			taskGateway, err := taskgateway.NewTaskGateway(logger, ctx, conf)
			if err != nil {
				return fmt.Errorf("failed to create TaskGateway: %w", err)
			}

			// start TaskDispatcher
			g.Go(func() error {
				err = taskDispatcher.Start()
				if err != nil {
					return fmt.Errorf("failed to start TaskDispatcher: %w", err)
				}
				<-taskDispatcher.Quit()
				return nil
			})

			// start TaskGateway
			g.Go(func() error {
				err = taskGateway.Start()
				if err != nil {
					return fmt.Errorf("failed to start TaskGateway: %w", err)
				}
				<-taskGateway.Quit()

				return nil
			})

			// wait for TaskDispatcher and TaskGateway to quit
			if err := g.Wait(); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&configFile, "config", "", "config file")
	return cmd
}
