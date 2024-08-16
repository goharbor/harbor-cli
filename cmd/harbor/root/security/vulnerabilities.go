package security

import (
	"fmt"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

func listVulnerabilitiesCommand() *cobra.Command {
	var query string

	cmd := &cobra.Command{
		Use:   "vulnerabilities",
		Short: "List vulnerabilities in the system",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := api.ListVulnerabilities(query)
			if err != nil {
				return fmt.Errorf("error listing vulnerabilities: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Query condition for filtering vulnerabilities")

	return cmd
}
