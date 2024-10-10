package project

import (
	"github.com/spf13/cobra"
)

func Project() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "project",
		Short:   "Manage projects and assign resources to them",
		Long:    `Manage projects in Harbor`,
		Example: `  harbor project list`,
	}
	cmd.AddCommand(
		CreateProjectCommand(),
		DeleteProjectCommand(),
		ListProjectCommand(),
		ViewCommand(),
		LogsProjectCommmand(),
		SearchProjectCommand(),
	)

	return cmd
}
