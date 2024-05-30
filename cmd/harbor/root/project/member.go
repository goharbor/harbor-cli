package project

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project/member"
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
		member.ListMemberCommand(),
		member.CreateMemberCommand(),
		member.DeleteMemberCommand(),
		member.ViewMemberCommand(),
		member.UpdateMemberCommand(),
	)

	return cmd
}
