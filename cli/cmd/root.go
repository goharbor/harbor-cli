package cmd

import (
	"os"

	"github.com/goharbor/harbor-cli/cli/cmd/artifact"
	"github.com/goharbor/harbor-cli/cli/cmd/project"
	"github.com/goharbor/harbor-cli/cli/cmd/registry"
	"github.com/goharbor/harbor-cli/cli/cmd/repository"
	"github.com/goharbor/harbor-cli/cli/cmd/version"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
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
	cmd.AddCommand(project.GetProjectCommand())
	cmd.AddCommand(registry.GetRegistryCommand())
	return cmd
}

// newListCommand creates a new `harbor list` command
func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [COMMAND]",
		Short: "list project, registry, artifact etc.",
		Long:  `List project, registry, artifact`,
	}

	cmd.PersistentFlags().String(constants.CredentialNameOption, "", constants.CredentialNameHelp)
	cmd.AddCommand(project.ListProjectCommand())
	cmd.AddCommand(registry.ListRegistryCommand())
	cmd.AddCommand(artifact.ListArtifactCommand())
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
	cmd.AddCommand(project.CreateProjectCommand())
	cmd.AddCommand(registry.CreateRegistryCommand())
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
	cmd.AddCommand(project.DeleteProjectCommand())
	cmd.AddCommand(registry.DeleteRegistryCommand())
	return cmd
}

// newUpdateCommand creates a new `harbor update` command
func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [COMMAND]",
		Short: "update registry, etc.",
		Long:  `Update registry`,
	}

	cmd.PersistentFlags().String(constants.CredentialNameOption, "", constants.CredentialNameHelp)
	cmd.AddCommand(registry.UpdateRegistryCommand())
	cmd.AddCommand(repository.UpdateRepositoryCommand())
	return cmd
}


var RootCmd = &cobra.Command{
	Use:   "harbor",
	Short: "Official Harbor CLI",
	Long:  `A Multi-Featured official Harbor CLI which can interact with the Harbor API`,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(version.VersionCommand(),
		LoginCommand(),
		newGetCommand(),
		newListCommand(),
		newCreateCommand(),
		newDeleteCommand(),
		newUpdateCommand())

}