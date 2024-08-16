package security

import "github.com/spf13/cobra"

func Security() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security",
		Short: "Manage security-related operations in Harbor",
	}

	cmd.AddCommand(
		getSecuritySummaryCommand(),
		listVulnerabilitiesCommand(),
	)

	return cmd
}