package member

import (
	"github.com/spf13/cobra"
)

func Member() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "member",
		Short:   "Manage member and assign resources to them",
		Long:    `Manage members in Harbor`,
		Example: `  harbor member list`,
	}
	cmd.AddCommand(
		ListMemberCommand(),
		CreateMemberCommand(),
    DeleteMemberCommand(),
	)

	return cmd
}
