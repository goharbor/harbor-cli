package retention

import (
	"github.com/spf13/cobra"
)

func Retention() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retention",
		Short: "Manage tag retention policies in the project",
		Long: `Manage tag retention policies in the project in Harbor.
		
The 'retention' command allows users to create, list, and delete tag retention rules 
within a project. Tag retention policies help in managing and controlling the lifecycle 
of tags by defining rules for automatic cleanup and retention.

A user can only create up to 15 tag retention rules per project.`,
		Example: `  harbor tag retention create    # Create a new tag retention policy
  harbor tag retention list      # List all tag retention rules in the project
  harbor tag retention delete    # Delete a specific tag retention policy`,
	}

	cmd.AddCommand(
		CreateRetentionCommand(),
		ListRetentionRulesCommand(),
		DeleteRetentionPolicyCommand(),
	)
	return cmd
}
