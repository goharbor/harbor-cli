package project

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DeleteProjectCommand creates a new `harbor delete project` command
func DeleteProjectCommand() *cobra.Command {
	var forceDelete bool

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete project by name or id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				err = api.DeleteProject(args[0], forceDelete)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				err = api.DeleteProject(projectName, forceDelete)
			}
			if err != nil {
				log.Errorf("failed to delete project: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&forceDelete, "force", false, "Deletes all repositories and artifacts within the project")

	return cmd
}
