package repository

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/repository/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListRepositoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list repositories within a project",
		Example: `  harbor repo list <project_name>`,
		Long:    `Get information of all repositories in a project`,
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var repos repository.ListRepositoriesOK
			var projectName string

			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName = prompt.GetProjectNameFromUser()
			}

			repos, err = api.ListRepository(projectName)

			if err != nil {
				log.Errorf("failed to list repositories: %v", err)
				return
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(repos, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.ListRepositories(repos.Payload)
			}
		},
	}

	return cmd
}
