package labels

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RemoveLabelsCmd() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:     "remove",
		Short:   "Remove a label of an artifact",
		Example: `harbor artifact labels remove <project>/<repository>/<reference> <labelName|labelID>`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var projectName, repoName, reference string
			var labelID int64

			if len(args) > 0 {
				projectName, repoName, reference = utils.ParseProjectRepoReference(args[0])
				labelID, err = api.GetLabelIdByName(args[1])
				if err != nil {
					logrus.Errorf("Failed to get this lable: %s", args[1])
					return
				}
			} else {
				projectName = prompt.GetProjectNameFromUser()
				repoName = prompt.GetRepoNameFromUser(projectName)
				reference = prompt.GetReferenceFromUser(repoName, projectName)
				labelID = prompt.GetLabelIdFromUser(opts)
			}
			err = api.RemoveLabel(projectName, repoName, reference, labelID)

			if err != nil {
				logrus.Errorf("Failed to remove label %s/%s@%s", projectName, repoName, reference)
				return
			}
		},
	}

	return cmd
}
