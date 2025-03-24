package retention

import (
	"github.com/spf13/cobra"
)

func Retention() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "retention",
		Short:   "Manage retention rule in the project",
		Long:    `Manage retention rules in the project in Harbor`,
		Example: `harbor retention create`,
	}
	cmd.AddCommand(
		CreateRetentionCommand(),
		ListExecutionRetentionCommand(),
		DeleteRetentionCommand(),
	)

	return cmd
}