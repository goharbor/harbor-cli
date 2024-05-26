package artifact

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func InfoArtifactCommmand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "info",
		Short:   "Get info of an artifact",
		Long:    `Get info of an artifact`,
		Example: `harbor artifact info <project>/<repository>/<reference>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				err = api.RunInfoArtifact(projectName, repoName, reference)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				reference := utils.GetReferenceFromUser(repoName, projectName)
				err = api.RunInfoArtifact(projectName, repoName, reference)
			}

			if err != nil {
				log.Errorf("failed to get info of an artifact: %v", err)
			}

		},
	}

	return cmd
}
