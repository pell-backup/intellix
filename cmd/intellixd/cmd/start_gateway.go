package cmd

import (
	"fmt"
	"intellix/gateway/dispatcher"
	"intellix/gateway/submitter"
	"intellix/gateway/types"

	dvslog "github.com/0xPellNetwork/pelldvs-libs/log"
	pelldvscfg "github.com/0xPellNetwork/pelldvs/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func taskGatewayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-task-gateway",
		Short: "Start the Submitter service",
		Long: "Start the Submitter service, Example:\n" +
			"intellixd start-task-gateway --config=config.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			//	serverCtx := server.GetServerContextFromCmd(cmd)
			logger := dvslog.NewNopLogger()

			home := getConfigHome()
			if configFile == "" {
				configFile = home + "/config/gateway.config.json"
			}

			conf, err := types.LoadConfig(configFile)
			if err != nil {
				return fmt.Errorf("failed to load Submitter configuration: %w", err)
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

			// create Dispatcher
			dispatcher, err := dispatcher.NewDispatcher(logger, pellDVSConf, conf)
			if err != nil {
				return fmt.Errorf("failed to create Dispatcher: %w", err)
			}

			// create Submitter
			submitter, err := submitter.NewSubmitter(logger, ctx, conf)
			if err != nil {
				return fmt.Errorf("failed to create Submitter: %w", err)
			}

			// start Dispatcher
			g.Go(func() error {
				err = dispatcher.Start()
				if err != nil {
					return fmt.Errorf("failed to start Dispatcher: %w", err)
				}
				<-dispatcher.Quit()
				return nil
			})

			// start Submitter
			g.Go(func() error {
				err = submitter.Start()
				if err != nil {
					return fmt.Errorf("failed to start Submitter: %w", err)
				}
				<-submitter.Quit()

				return nil
			})

			// wait for Dispatcher and Submitter to quit
			if err := g.Wait(); err != nil {
				return err
			}

			return nil
		},
	}
	cmd.Flags().StringVar(&configFile, "config", "", "config file")
	return cmd
}
