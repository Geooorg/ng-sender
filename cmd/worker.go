package cmd

import (
	"flag"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/worker"
	"log"
	wf "ng-sender/pkg/workflow"
)

// starts a Temporal worker
var workerCmd = &cobra.Command{
	Use:   "temporal-worker",
	Short: "Runs a worker for temporal workflows",
	Run: func(cmd *cobra.Command, args []string) {
		flag.Parse()

		cfg := &config{}
		if err := viper.Unmarshal(cfg); err != nil {
			log.Fatal(err)
		}

		temporalClient, err := getTemporalClient(cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer temporalClient.Close()

		natsClient, err := getNatsClient(cfg)
		if err != nil {
			fmt.Printf("unable to create connection %s\n", err)
			fmt.Printf("nats config: %v\n", cfg.Nats)
			return
		}
		defer natsClient.Close()

		w := worker.New(temporalClient, "warningMessages", worker.Options{})
		defer w.Stop()

		activities := &wf.WarningMessageActivities{
			NatsClient:   natsClient,
			TopicsConfig: cfg.Nats.TopicsConfig,
		}

		w.RegisterActivity(activities)
		w.RegisterWorkflow(wf.SendToReceiversWF)

		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalln("Unable to start worker", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
