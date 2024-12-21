package config

import "github.com/spf13/cobra"

func Config() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage the config of the Harbor Cli",
		Long:  `Manage repositories in Harbor config`,
	}
	cmd.AddCommand(
		ListConfigCommand(),
		GetConfigItemCommand(),
		SetConfigItemCommand(),
		DeleteConfigItemCommand(),
	)

	return cmd

}
