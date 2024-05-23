package artifact

import (
	"context"
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/utils"
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

			if len(args) > 0 {
				projectName, repoName := utils.ParseProjectRepo(args[0])
				err = runListArtifact(projectName, repoName)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				err = runListArtifact(projectName, repoName)
			}

			if err != nil {
				log.Errorf("failed to list artifacts: %v", err)
			}

		},
	}

	return cmd
}

func runListArtifact(projectName, repoName string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	response, err := client.Artifact.ListArtifacts(ctx, &artifact.ListArtifactsParams{ProjectName: projectName, RepositoryName: repoName})

	if err != nil {
		return err
	}

	fmt.Println(response.Payload)

	return nil

}
