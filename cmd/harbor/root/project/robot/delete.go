package robot

import (
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// to-do improve DeleteRobotCommand and multi delete
func DeleteRobotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [robotID]",
		Short: "delete robot by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				robotID, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Errorf("failed to parse robot ID: %v", err)
				}
				err = api.DeleteRobot(robotID)
				if err != nil {
					log.Errorf("failed to Delete robots")
				}
			}
		},
	}

	return cmd
}
