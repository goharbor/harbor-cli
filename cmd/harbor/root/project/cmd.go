package project

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project/metadata"
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
		metadata.Metadata(),
	)

	return cmd
}
