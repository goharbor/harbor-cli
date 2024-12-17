package user

import (
	"github.com/spf13/cobra"
)

func User() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user",
		Short:   "Manage users",
		Long:    `Manage users in Harbor`,
		Example: `  harbor user list`,
	}

	cmd.AddCommand(
		UserListCmd(),
		UserCreateCmd(),
		UserDeleteCmd(),
		ElevateUserCmd(),
		UserResetCmd(),
	)

	return cmd
}
