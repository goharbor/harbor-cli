package member

import (
	"github.com/spf13/cobra"
)

func Member() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "member",
		Short:   `Manage members in a Project`,
		Long:    "Manage members and assign roles to them",
		Example: `  harbor member list`,
	}
	cmd.AddCommand(
		ListMemberCommand(),
		CreateMemberCommand(),
		DeleteMemberCommand(),
		ViewMemberCommand(),
    UpdateMemberCommand(),
	)

	return cmd
}
