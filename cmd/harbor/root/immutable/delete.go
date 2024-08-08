package immutable

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
	"github.com/goharbor/harbor-cli/pkg/prompt"
)

func DeleteImmutableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete immutable rule",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var immutableId int64
			if len(args) > 0 {
				immutableId = prompt.GetImmutableTagRule(args[0])
				err = api.DeleteImmutable(args[0],immutableId)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				immutableId = prompt.GetImmutableTagRule(projectName)
				err = api.DeleteImmutable(projectName,immutableId)
			}
			if err != nil {
				log.Errorf("failed to delete immutable rule: %v", err)
			}
		},
	}

	return cmd
}