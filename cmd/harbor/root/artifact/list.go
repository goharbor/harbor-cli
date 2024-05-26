package artifact

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ListArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list artifacts within a repository",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				projectName, repoName := utils.ParseProjectRepo(args[0])
				err = api.RunListArtifact(projectName, repoName)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				err = api.RunListArtifact(projectName, repoName)
			}

			if err != nil {
				log.Errorf("failed to list artifacts: %v", err)
			}

		},
	}

	return cmd
}
