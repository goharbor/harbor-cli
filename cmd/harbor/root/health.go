package root

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
	"github.com/goharbor/harbor-cli/pkg/views/health"
)

func HealthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Get the health status of Harbor components",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := api.GetHealth()
			if err != nil {
				return err
			}
			health.PrintHealthStatus(status)
			return nil
		},
	}

	return cmd
}
