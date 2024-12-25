package securityhub

import (
	"github.com/spf13/cobra"
)

func SecurityHub() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security-hub",
		Short: "Security Hub for managing security vulnerability",
		Long:  "Security Hub provides tools to manage, monitor, and remediate security issues in repositories",
	}
	cmd.AddCommand(
		ListVulnerabilityCommand(),
	)

	return cmd
}
