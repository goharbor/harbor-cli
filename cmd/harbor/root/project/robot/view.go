package robot

import (
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/robot"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/robot/list"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// ViewCommand creates a new `harbor project robot view` command
func ViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [robotID]",
		Short: "get robot by id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var robot *robot.GetRobotByIDOK

			if len(args) == 1 {
				robotID, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Errorf("failed to parse robot ID: %v", err)
				}
				robot, err = api.GetRobot(robotID)
				if err != nil {
					log.Errorf("failed to List robots")
				}
			}
			robots := []*models.Robot{robot.Payload}
			list.ListRobots(robots)
		},
	}

	return cmd
}
