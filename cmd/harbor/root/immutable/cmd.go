package immutable

import (
	"github.com/spf13/cobra"
)

func Immutable() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "immutable",
		Short:   "Manage Immutability rules in the project",
		Long:    `Manage Immutability rules in the project in Harbor`,
		Example: `harbor immutable create`,
	}
	cmd.AddCommand(
		CreateImmutableCommand(),
		ListImmutableCommand(),
		DeleteImmutableCommand(),
	)

	return cmd
}