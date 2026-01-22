package gc

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RunGCCommand() *cobra.Command {
	var dryRun, deleteUntagged bool

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run Garbage Collection manually",
		Run: func(cmd *cobra.Command, args []string) {

			scheduleObj := models.ScheduleObj{
				Type: "Manual",
			}

			params := map[string]interface{}{
				"dry_run":         dryRun,
				"delete_untagged": deleteUntagged,
			}

			scheduleBody := &models.Schedule{
				Schedule:   &scheduleObj,
				Parameters: params,
			}

			err := api.CreateGCSchedule(scheduleBody)
			if err != nil {
				logrus.Fatalf("Failed to start GC: %v", err)
			}
			logrus.Info("GC started successfully")
		},
	}

	cmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "Simulate GC without deleting artifacts")
	cmd.Flags().BoolVarP(&deleteUntagged, "delete-untagged", "", true, "Delete untagged artifacts")

	return cmd
}
