package scan_all

import (
	"fmt"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

func GetScanAllMetricsCommand() *cobra.Command {
	var scheduled bool

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Get the metrics of the latest scan all process",
		RunE: func(cmd *cobra.Command, args []string) error {
			metrics, err := api.GetScanAllMetrics(scheduled)
			if err != nil {
				return err
			}

			fmt.Println("Total: ", metrics.Total)
			fmt.Println("Ongoing: ", metrics.Ongoing)
			fmt.Println("Completed: ", metrics.Completed)
			fmt.Println("Trigger: ", metrics.Trigger)
			fmt.Println("Metrics: ", metrics.Metrics)

			return nil
		},
	}

	flags := cmd.Flags()
	// latest scheduled metrics is deprecated in the API
	flags.BoolVarP(&scheduled, "scheduled", "s", false, "Get the metrics of the latest scheduled scan all process")

	return cmd
}
