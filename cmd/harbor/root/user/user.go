package user

import "github.com/spf13/cobra"

func UserCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "user",
		Short: "manage users",
	}

	cmd.AddCommand(UserListCmd())
	cmd.AddCommand(UserCreateCmd())

	return cmd
}
