package root

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/health"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func HealthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Get the health status of Harbor components",
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				log.Fatalf("Error: accepts 0 arg(s), received %d: %v", len(args), args)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := api.GetHealth()
			if err != nil {
				return err
			}
			health.PrintHealthStatus(status)
			return nil
		},
		Example: `  # Get the health status of Harbor components`,
	}

	return cmd

}
