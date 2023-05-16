package cmd

import (
	"github.com/akshatdalton/harbor-cli/cmd/constants"
	"github.com/akshatdalton/harbor-cli/cmd/login"
	"github.com/akshatdalton/harbor-cli/cmd/project"
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

	cmd.PersistentFlags().String(constants.CredentialNameOption, "", constants.CredentialNameHelp)
	cmd.AddCommand(project.NewGetProjectCommand())
	cmd.AddCommand(registry.NewGetRegistryCommand())
	return cmd
}

// newListCommand creates a new `harbor list` command
func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [COMMAND]",
		Short: "list project, registry, etc.",
		Long:  `List project, registry`,
	}

	cmd.PersistentFlags().String(constants.CredentialNameOption, "", constants.CredentialNameHelp)
	cmd.AddCommand(project.NewListProjectCommand())
	cmd.AddCommand(registry.NewListRegistryCommand())
	return cmd
}

// newCreateCommand creates a new `harbor create` command
func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [COMMAND]",
		Short: "create project, registry, etc.",
		Long:  `Create project, registry`,
	}

	cmd.PersistentFlags().String(constants.CredentialNameOption, "", constants.CredentialNameHelp)
	cmd.AddCommand(project.NewCreateProjectCommand())
	cmd.AddCommand(registry.NewCreateRegistryCommand())
	return cmd
}

// newDeleteCommand creates a new `harbor delete` command
func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [COMMAND]",
		Short: "delete project, registry, etc.",
		Long:  `Delete project, registry`,
	}

	cmd.PersistentFlags().String(constants.CredentialNameOption, "", constants.CredentialNameHelp)
	cmd.AddCommand(project.NewDeleteProjectCommand())
	cmd.AddCommand(registry.NewDeleteRegistryCommand())
	return cmd
}

func addCommands(cmd *cobra.Command) {
	cmd.AddCommand(login.NewLoginCommand())
	cmd.AddCommand(newGetCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newCreateCommand())
	cmd.AddCommand(newDeleteCommand())
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
