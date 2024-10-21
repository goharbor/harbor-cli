package artifact

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	artifactViews "github.com/goharbor/harbor-cli/pkg/views/artifact/list"
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
			var resp artifact.ListArtifactsOK

			if len(args) > 0 {
				projectName, repoName := utils.ParseProjectRepo(args[0])
				resp, err = api.ListArtifact(projectName, repoName)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				resp, err = api.ListArtifact(projectName, repoName)
			}

			if err != nil {
				log.Errorf("failed to list artifacts: %v", err)
			}
			artifactViews.ListArtifacts(resp.Payload)
		},
	}

	return cmd
}
