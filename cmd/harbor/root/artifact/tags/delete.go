package tags

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeleteTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a tag of an artifact",
		Example: `harbor artifact tags delete <project>/<repository>/<reference> <tag>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				tag := args[1]
				err = runDeleteTags(projectName, repoName, reference, tag)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				reference := utils.GetReferenceFromUser(repoName, projectName)
				// tag := utils.GetTagFromUseer(repoName, projectName)
				tag := "sa"
				err = runDeleteTags(projectName, repoName, reference, tag)
			}
			if err != nil {
				log.Errorf("failed to delete a tag of an artifact: %v", err)
			}
		},
	}

	return cmd
}

func runDeleteTags(projectName, repoName, reference, tag string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	_, err := client.Artifact.DeleteTag(ctx, &artifact.DeleteTagParams{ProjectName: projectName, RepositoryName: repoName, Reference: reference, TagName: tag})
	if err != nil {
		return err
	}
	log.Infof("Tag deleted successfully")
	return nil
}
