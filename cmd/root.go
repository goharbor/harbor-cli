package cmd

import (
	"github.com/akshatdalton/harbor-cli/cmd/login"
	"github.com/akshatdalton/harbor-cli/cmd/registry"
	"github.com/spf13/cobra"
)

// newGetCommand creates a new `harbor get` command
func newGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [COMMAND]",
		Short: "get project, registry, etc.",
		Long:  `Get project, registry`,
	}

	cmd.AddCommand(registry.NewGetRegistryCommand())
	return cmd
}

func addCommands(cmd *cobra.Command) {
	cmd.AddCommand(login.NewLoginCommand())
	cmd.AddCommand(newGetCommand())
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
