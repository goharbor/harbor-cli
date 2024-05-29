package artifact

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
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
				err = api.InfoArtifact(projectName, repoName, reference)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				reference := prompt.GetReferenceFromUser(repoName, projectName)
				err = api.InfoArtifact(projectName, repoName, reference)
			}

			if err != nil {
				log.Errorf("failed to get info of an artifact: %v", err)
			}

		},
	}

	return cmd
}
