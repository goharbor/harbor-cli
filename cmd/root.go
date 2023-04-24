package cmd

import (
	"github.com/akshatdalton/harbor-cli/cmd/login"
	"github.com/spf13/cobra"
)

func addCommands(cmd *cobra.Command) {
	cmd.AddCommand(login.NewLoginCommand())
}

// CreateHarborCLI creates a new Harbor CLI
func CreateHarborCLI() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "harbor",
		Short: "Official Harbor CLI",
	}

	addCommands(cmd)
	return cmd
}
