package security

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/security/summary"
	"github.com/spf13/cobra"
)

func getSecuritySummaryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Get the security summary of the system",
		RunE: func(cmd *cobra.Command, args []string) error {
			response,err := api.GetSecuritySummary()
			if err != nil {
				return fmt.Errorf("error getting security summary: %w", err)
			}
			summary.SecuritySummary(response.Payload)
			return nil
		},
	}

	return cmd
}