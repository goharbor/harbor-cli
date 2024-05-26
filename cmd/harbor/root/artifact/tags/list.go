package tags

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListTagsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tags of an artifact",
		Example: `harbor artifact tags list <project>/<repository>/<reference>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				err = runListTags(projectName, repoName, reference)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				reference := utils.GetReferenceFromUser(repoName, projectName)
				err = runListTags(projectName, repoName, reference)
			}
			if err != nil {
				log.Errorf("failed to list tags of an artifact: %v", err)
			}

		},
	}

	return cmd
}

func runListTags(projectName, repoName, reference string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	resp, err := client.Artifact.ListTags(ctx, &artifact.ListTagsParams{ProjectName: projectName, RepositoryName: repoName, Reference: reference})
	if err != nil {
		return err
	}

	log.Info(resp.Payload)
	return nil
}
