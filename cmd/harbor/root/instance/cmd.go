package instance

import "github.com/spf13/cobra"

func Instance() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "Manage instance in Harbor",
	}
	cmd.AddCommand(
		CreateInstanceCommand(),
		DeleteInstanceCommand(),
		ListInstanceCommand(),
	)
	return cmd
}