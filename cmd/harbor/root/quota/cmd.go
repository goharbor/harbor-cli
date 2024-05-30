package quota

import (
	"github.com/spf13/cobra"
)

func Quota() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "quota",
		Short:   "Manage quotas",
		Long:    `Manage quotas`,
		Example: `  harbor quota list`,
	}
	cmd.AddCommand(
		ListQuotaCommand(),
	)

	return cmd
}
