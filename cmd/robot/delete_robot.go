package robot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/akshatdalton/harbor-cli/cmd/constants"
	"github.com/akshatdalton/harbor-cli/cmd/utils"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/robot"
	"github.com/spf13/cobra"
)

type deleteRobotOptions struct {
	id int64
}

// NewDeleteRobotCommand creates a new `harbor delete robot` command
func NewDeleteRobotCommand() *cobra.Command {
	var opts deleteRobotOptions

	cmd := &cobra.Command{
		Use:   "robot [ID]",
		Short: "delete robot by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid argument: %s. Expected an integer.\n", args[0])
				return err
			}
			opts.id = int64(id)

			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return runDeleteRobot(opts, credentialName)
		},
	}

	return cmd
}

func runDeleteRobot(opts deleteRobotOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Robot.DeleteRobot(ctx, &robot.DeleteRobotParams{RobotID: opts.id})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
