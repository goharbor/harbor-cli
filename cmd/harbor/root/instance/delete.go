package instance

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

func DeleteInstanceCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete instance by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				instanceName := args[0]
				err = api.DeleteInstance(instanceName)
			} else {
				instanceName := prompt.GetInstanceFromUser()
				err = api.DeleteInstance(instanceName)
			}
			if err != nil {
				log.Errorf("failed to delete instance: %v", err)
			}
		},
	}

	return cmd
}