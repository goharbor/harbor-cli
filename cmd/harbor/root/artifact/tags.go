package artifact

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ArtifactTagsCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "tags",
		Short:   "Manage tags of an artifact",
		Example: ` harbor artifact tags list <project>/<repository>/<reference>`,
	}

	cmd.AddCommand(
		ListTagsCommand(),
		DeleteTagsCmd(),
	)

	return cmd
}

func ListTagsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tags of an artifact",
		Example: `harbor artifact tags list <project>/<repository>/<reference>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				err = api.ListTags(projectName, repoName, reference)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				reference := utils.GetReferenceFromUser(repoName, projectName)
				err = api.ListTags(projectName, repoName, reference)
			}
			if err != nil {
				log.Errorf("failed to list tags of an artifact: %v", err)
			}

		},
	}

	return cmd
}

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
				err = api.DeleteTag(projectName, repoName, reference, tag)
			} else {
				projectName := utils.GetProjectNameFromUser()
				repoName := utils.GetRepoNameFromUser(projectName)
				reference := utils.GetReferenceFromUser(repoName, projectName)
				tag := utils.GetTagFromUser(repoName, projectName, reference)
				err = api.DeleteTag(projectName, repoName, reference, tag)
			}
			if err != nil {
				log.Errorf("failed to delete a tag of an artifact: %v", err)
			}
		},
	}

	return cmd
}
