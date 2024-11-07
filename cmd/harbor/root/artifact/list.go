package artifact

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	artifactViews "github.com/goharbor/harbor-cli/pkg/views/artifact/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list artifacts within a repository",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var artifacts artifact.ListArtifactsOK
			var projectName, repoName string

			if len(args) > 0 {
				projectName, repoName = utils.ParseProjectRepo(args[0])
			} else {
				projectName = prompt.GetProjectNameFromUser()
				repoName = prompt.GetRepoNameFromUser(projectName)
			}

			artifacts, err = api.ListArtifact(projectName, repoName)

			if err != nil {
				log.Errorf("failed to list artifacts: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(artifacts, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				artifactViews.ListArtifacts(artifacts.Payload)
			}
		},
	}

	return cmd
}
