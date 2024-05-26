package repository

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ListRepositoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list repositories within a project",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				err = api.ListRepository(args[0])
			} else {
				projectName := utils.GetProjectNameFromUser()
				err = api.ListRepository(projectName)
			}
			if err != nil {
				log.Errorf("failed to list repositories: %v", err)
			}
		},
	}

	return cmd
}
