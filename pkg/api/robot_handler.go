package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/permissions"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/robot"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	"github.com/goharbor/harbor-cli/pkg/views/robot/update"
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

func CreateRobot(opts create.CreateView) (*robot.CreateRobotCreated, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	// Create a slice to store converted permissions
	permissions := opts.Permissions
	convertedPerms := make([]*models.RobotPermission, 0, len(permissions))

	project := "project"
	// Loop through original permissions and convert them
	for _, perm := range permissions {
		convertedPerm := &models.RobotPermission{
			Access:    perm.Access,
			Kind:      project,
			Namespace: opts.ProjectName,
		}
		convertedPerms = append(convertedPerms, convertedPerm)
	}
	response, err := client.Robot.CreateRobot(
		ctx,
		&robot.CreateRobotParams{
			Robot: &models.RobotCreate{
				Description: opts.Description,
				Disable:     false,
				Duration:    opts.Duration,
				Level:       opts.Level,
				Name:        opts.Name,
				Permissions: convertedPerms,
			},
		},
	)
	if err != nil {
		return nil, err
	}

	log.Info("robot created successfully.")
	return response, nil
}

// update robot with robotID
func UpdateRobot(opts *update.UpdateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		log.Errorf("Error: %v", err)
		return err
	}

	log.Println(opts)

	// Create a slice to store converted permissions
	permissions := opts.Permissions
	convertedPerms := make([]*models.RobotPermission, 0, len(permissions))

	kind := "project"
	// Loop through original permissions and convert them
	for _, perm := range permissions {
		convertedPerm := &models.RobotPermission{
			Access:    perm.Access,
			Kind:      kind,
			Namespace: opts.Permissions[0].Namespace,
		}
		convertedPerms = append(convertedPerms, convertedPerm)
	}
	_, err = client.Robot.UpdateRobot(
		ctx,
		&robot.UpdateRobotParams{
			Robot: &models.Robot{
				Description: opts.Description,
				Duration:    opts.Duration,
				Editable:    opts.Editable,
				Disable:     opts.Disable,
				ID:          opts.ID,
				Level:       opts.Level,
				Name:        opts.Name,
				Permissions: convertedPerms,
			},
			RobotID: opts.ID,
		},
	)
	if err != nil {
		log.Errorf("Error in updating Robot: %v", err)
		return err
	}

	log.Info("robot updated successfully.")
	return nil
}

func GetPermissions() (*permissions.GetPermissionsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	response, err := client.Permissions.GetPermissions(
		ctx,
		&permissions.GetPermissionsParams{},
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func RefreshSecret(secret string, robotID int64) (*robot.RefreshSecOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	robotSec := &models.RobotSec{
		Secret: secret,
	}

	response, err := client.Robot.RefreshSec(ctx, &robot.RefreshSecParams{
		RobotSec: robotSec,
		RobotID:  robotID,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
