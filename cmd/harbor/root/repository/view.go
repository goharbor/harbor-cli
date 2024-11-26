package repository

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/repository/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RepoViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "view",
		Short:   "Get repository information",
		Example: `  harbor repo view <project_name>/<repo_name>`,
		Long:    `Get information of a particular repository in a project`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var projectName, repoName string
			var repo *repository.GetRepositoryOK

			if len(args) > 0 {
				projectName, repoName = utils.ParseProjectRepo(args[0])
			} else {
				projectName = prompt.GetProjectNameFromUser()
				repoName = prompt.GetRepoNameFromUser(projectName)
			}

			repo, err = api.RepoView(projectName, repoName)

			if err != nil {
				log.Errorf("failed to get repository information: %v", err)
				return
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(repo, FormatFlag)
				if err != nil {
					log.Error(err)
					return
				}
			} else {
				view.ViewRepository(repo.Payload)
			}

		},
	}

	return cmd
}
