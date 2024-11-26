package artifact

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/artifact/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ViewArtifactCommmand() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "view",
		Short:   "Get information of an artifact",
		Long:    `Get information of an artifact`,
		Example: `harbor artifact view <project>/<repository>/<reference>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var projectName, repoName, reference string
			var artifact *artifact.GetArtifactOK

			if len(args) > 0 {
				projectName, repoName, reference = utils.ParseProjectRepoReference(args[0])
			} else {
				projectName = prompt.GetProjectNameFromUser()
				repoName = prompt.GetRepoNameFromUser(projectName)
				reference = prompt.GetReferenceFromUser(repoName, projectName)
			}

			artifact, err = api.ViewArtifact(projectName, repoName, reference)

			if err != nil {
				log.Errorf("failed to get info of an artifact: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(artifact, FormatFlag)
				if err != nil {
					log.Error(err)
					return
				}
			} else {
				view.ViewArtifact(artifact.Payload)
			}

		},
	}

	return cmd
}
