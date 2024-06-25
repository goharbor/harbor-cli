package project

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// GetProjectCommand creates a new `harbor get project` command
func ViewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "view [NAME|ID]",
		Short: "get project by name or id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				err = api.GetProject(args[0])
			} else {
				projectName := prompt.GetProjectNameFromUser()
				err = api.GetProject(projectName)
			}

			if err != nil {
				log.Errorf("failed to get project: %v", err)
			}

		},
	}

	return cmd
}
