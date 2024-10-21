package artifact

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/artifact/tags/create"
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
		ListTagsCmd(),
		DeleteTagsCmd(),
		CreateTagsCmd(),
	)

	return cmd
}

func CreateTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a tag of an artifact",
		Example: `harbor artifact tags create <project>/<repository>/<reference> <tag>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				tag := args[1]
				err = api.CreateTag(projectName, repoName, reference, tag)
			} else {
				var tagName string
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				reference := prompt.GetReferenceFromUser(repoName, projectName)
				create.CreateTagView(&tagName)
				err = api.CreateTag(projectName, repoName, reference, tagName)
			}
			if err != nil {
				log.Errorf("failed to create tag: %v", err)
			}
		},
	}

	return cmd
}

func ListTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tags of an artifact",
		Example: `harbor artifact tags list <project>/<repository>/<reference>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			var resp artifact.ListTagsOK
			if len(args) > 0 {
				projectName, repoName, reference := utils.ParseProjectRepoReference(args[0])
				resp, err = api.ListTags(projectName, repoName, reference)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				reference := prompt.GetReferenceFromUser(repoName, projectName)
				resp, err = api.ListTags(projectName, repoName, reference)
			}
			if err != nil {
				log.Errorf("failed to list tags: %v", err)
			}

			log.Info(resp.Payload)
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
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				reference := prompt.GetReferenceFromUser(repoName, projectName)
				tag := prompt.GetTagFromUser(repoName, projectName, reference)
				err = api.DeleteTag(projectName, repoName, reference, tag)
			}
			if err != nil {
				log.Errorf("failed to delete tag: %v", err)
			}
		},
	}

	return cmd
}
