package repository

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RepoInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "Get repository information",
		Example: `  harbor repo info <project_name>/<repo_name>`,
		Long:    `Get information of a particular repository in a project`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				projectName, repoName := utils.ParseProjectRepo(args[0])
				err = api.RepoInfo(projectName, repoName)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				err = api.RepoInfo(projectName, repoName)
			}
			if err != nil {
				log.Errorf("failed to get repository information: %v", err)
			}

		},
	}

	return cmd
}
