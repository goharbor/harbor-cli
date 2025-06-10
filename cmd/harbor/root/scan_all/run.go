package scan_all

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

func RunScanAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Scan all artifacts now",
		RunE: func(cmd *cobra.Command, args []string) error {
			return api.UpdateScanAllSchedule(models.ScheduleObj{Type: "Manual"})
		},
	}

	return cmd
}
