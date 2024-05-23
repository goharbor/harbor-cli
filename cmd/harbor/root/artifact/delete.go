package artifact

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeleteArtifactCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete an artifact",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				err = runDeleteArtifact(projectName, repoName, reference)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				reference := utils.GetReferenceFromUser(repoName, projectName)
				err = runDeleteArtifact(projectName, repoName, reference)
			}

			if err != nil {
				log.Errorf("failed to delete an artifact: %v", err)
			}
		},
	}

	return cmd
}

func runDeleteArtifact(projectName, repoName, reference string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	_, err := client.Artifact.DeleteArtifact(ctx, &artifact.DeleteArtifactParams{ProjectName: projectName, RepositoryName: repoName, Reference: reference})

	if err != nil {
		return err
	}

	log.Infof("Artifact deleted successfully")

	return nil
}
