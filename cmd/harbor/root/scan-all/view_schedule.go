package scan_all

import (
	"fmt"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

// This command does not work because the API does not return the response body
// API: https://demo.goharbor.io/devcenter-api-2.0
func ViewScanAllScheduleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view-schedule",
		Short:   "View the scan all schedule",
		Aliases: []string{"vs"},
		RunE: func(cmd *cobra.Command, args []string) error {
			schedule, err := api.GetScanAllSchedule()
			if err != nil {
				return err
			}
			fmt.Println("Current cron: ", schedule.Cron)
			fmt.Println("Current next scan time: ", schedule.NextScheduledTime)
			fmt.Println("Current scan all type: ", schedule.Type)
			return nil
		},
	}

	return cmd
}
