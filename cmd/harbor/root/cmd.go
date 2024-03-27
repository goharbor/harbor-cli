package root

import (
	"fmt"

	"github.com/goharbor/harbor-cli/cmd/harbor/internal/version"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/registry"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/user"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/spf13/cobra"
)

// versionCommand creates a new `harbor version` command
func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "get Harbor CLI version",
		Long:  `Get Harbor CLI version, git commit, go version, build time, release channel, os/arch, etc.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version:      %s\n", version.Version)
			fmt.Printf("Go version:   %s\n", version.GoVersion)
			fmt.Printf("Git commit:   %s\n", version.GitCommit)
			fmt.Printf("Built:        %s\n", version.BuildTime)
			fmt.Printf("OS/Arch:      %s\n", version.System)
		},
	}
}

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
		Short: "list project, registry, etc.",
		Long:  `List project, registry`,
	}

	cmd.PersistentFlags().String(constants.CredentialNameOption, "", constants.CredentialNameHelp)
	cmd.AddCommand(project.ListProjectCommand())
	cmd.AddCommand(registry.ListRegistryCommand())
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
	return cmd
}

// CreateHarborCLI creates a new Harbor CLI
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "harbor [command]",
		Short: "Official Harbor CLI",
	}

	cmd.AddCommand(
		versionCommand(),
		LoginCommand(),
		newGetCommand(),
		newListCommand(),
		newCreateCommand(),
		newDeleteCommand(),
		newUpdateCommand(),
		user.UserCmd(),
	)
	return cmd
}
