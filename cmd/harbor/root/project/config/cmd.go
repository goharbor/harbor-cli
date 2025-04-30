package config

import "github.com/spf13/cobra"

func ProjectConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage project metadata",
	}
	cmd.AddCommand(
		AddConfigCommand(),
		DeleteConfigCommand(),
		ViewConfigCommand(),
		UpdateConfigCommand(),
		ListConfigCommand(),
	)

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "Use project ID instead of name")

	return cmd
}
