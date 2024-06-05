package project

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project/robot"
	"github.com/spf13/cobra"
)

func Robot() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "robot",
		Short:   "Manage robot accounts",
		Example: `  harbor project robot list`,
	}
	cmd.AddCommand(
		robot.ListRobotCommand(),
		robot.DeleteRobotCommand(),
		robot.ViewRobotCommand(),
		robot.CreateRobotCommand(),
		robot.UpdateRobotCommand(),
		robot.RefreshSecretCommand(),
	)

	return cmd
}
