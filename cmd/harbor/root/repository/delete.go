package repository

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RepoDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a repository",
		Example: `  harbor repository delete [project_name]/[repository_name]`,
		Long:    `Delete a repository within a project in Harbor`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				projectName, repoName := utils.ParseProjectRepo(args[0])
				err = api.RepoDelete(projectName, repoName)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				err = api.RepoDelete(projectName, repoName)
			}
			if err != nil {
				log.Errorf("failed to delete repository: %v", err)
			}
		},
	}
	return cmd
}
