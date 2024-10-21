package artifact

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
				err = api.DeleteArtifact(projectName, repoName, reference)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				repoName := prompt.GetRepoNameFromUser(projectName)
				reference := prompt.GetReferenceFromUser(repoName, projectName)
				err = api.DeleteArtifact(projectName, repoName, reference)
			}

			if err != nil {
				log.Errorf("failed to delete an artifact: %v", err)
			}
		},
	}

	return cmd
}
