package repository

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/repository/list"
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
			var resp repository.ListRepositoriesOK

			if len(args) > 0 {
				resp, err = api.ListRepository(args[0])
			} else {
				projectName := prompt.GetProjectNameFromUser()
				resp, err = api.ListRepository(projectName)
			}

			if err != nil {
				log.Errorf("failed to list repositories: %v", err)
			}

			list.ListRepositories(resp.Payload)

		},
	}

	return cmd
}
