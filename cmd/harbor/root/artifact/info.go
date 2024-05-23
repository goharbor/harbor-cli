package artifact

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
				err = runInfoArtifact(projectName, repoName, reference)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				reference := utils.GetReferenceFromUser(repoName, projectName)
				err = runInfoArtifact(projectName, repoName, reference)
			}

			if err != nil {
				log.Errorf("failed to get info of an artifact: %v", err)
			}

		},
	}

	return cmd
}

func runInfoArtifact(projectName, repoName, reference string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	response, err := client.Artifact.GetArtifact(ctx, &artifact.GetArtifactParams{ProjectName: projectName, RepositoryName: repoName, Reference: reference})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.Payload)

	return nil
}
