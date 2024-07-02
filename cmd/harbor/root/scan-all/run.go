package scan_all

import (
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

func RunScanAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Scan all artifacts now",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Random cron expression and random time need to be passed to the API, even though they are not used, otherwise it returns bad request
			randomCron := "0 * * * * *"
			randomTime := strfmt.DateTime{}
			return api.CreateScanAllSchedule(models.ScheduleObj{Type: "Manual", Cron: randomCron, NextScheduledTime: randomTime})
		},
	}

	return cmd
}
