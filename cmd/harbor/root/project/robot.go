package project

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project/robot"
	"github.com/spf13/cobra"
)

func Robot() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "robot",
		Short:   "Manage robot accounts",
		Example: `  harbor robot list`,
	}
	cmd.AddCommand(
		robot.ListRobotCommand(),
    robot.DeleteCommand(),
    robot.ViewCommand(),
    robot.CreateRobot(),
	)

	return cmd
}
