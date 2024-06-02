package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/robot"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func ListRobot(opts ListFlags) (*robot.ListRobotOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Robot.ListRobot(
		ctx,
		&robot.ListRobotParams{
			Page:     &opts.Page,
			PageSize: &opts.PageSize,
			Q:        &opts.Q,
			Sort:     &opts.Sort,
		},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetRobot(robotID int64) (*robot.GetRobotByIDOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Robot.GetRobotByID(ctx, &robot.GetRobotByIDParams{RobotID: robotID})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func DeleteRobot(robotID int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Robot.DeleteRobot(ctx, &robot.DeleteRobotParams{RobotID: robotID})
	if err != nil {
		return err
	}

	log.Info("robot deleted successfully")

	return nil
}
