package robot

import (
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// to-do improve DeleteRobotCommand and multi select & delete
func DeleteRobotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [robotID]",
		Short: "delete robot by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				robotID int64
				err     error
			)
			if len(args) == 1 {
				robotID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Fatalf("failed to parse robot ID: %v", err)
				}
			} else {
				projectID := prompt.GetProjectIDFromUser()
				robotID = prompt.GetRobotIDFromUser(projectID)
			}
			err = api.DeleteRobot(robotID)
			if err != nil {
				log.Fatalf("failed to Delete robots")
			}
		},
	}

	return cmd
}
