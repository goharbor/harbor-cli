package robot

import (
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/robot/update"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// to-do complete UpdateRobotCommand
func UpdateRobotCommand() *cobra.Command {
	var (
		robotID int64
		opts    update.UpdateView
		all     bool
	)

	cmd := &cobra.Command{
		Use:   "update [robotID]",
		Short: "update robot by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) == 1 {
				robotID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Errorf("failed to parse robot ID: %v", err)
				}

			} else {
				projectID := prompt.GetProjectIDFromUser()
				robotID = prompt.GetRobotIDFromUser(projectID)
			}

			robot, err := api.GetRobot(robotID)
			bot := robot.Payload

			opts = update.UpdateView{
				CreationTime: bot.CreationTime,
				Description:  bot.Description,
				Disable:      bot.Disable,
				Duration:     bot.Duration,
				Editable:     bot.Editable,
				ID:           bot.ID,
				Level:        bot.Level,
				Name:         bot.Name,
				Secret:       bot.Secret,
			}

			// declare empty permissions to hold permissions
			permissions := []models.Permission{}

			if all {
				perms, _ := api.GetPermissions()
				permission := perms.Payload.Project

				choices := []models.Permission{}
				for _, perm := range permission {
					choices = append(choices, *perm)
				}
				permissions = choices
			} else {
				permissions = prompt.GetRobotPermissionsFromUser()
			}

			// []Permission to []*Access
			var accesses []*models.Access
			for _, perm := range permissions {
				access := &models.Access{
					Action:   perm.Action,
					Resource: perm.Resource,
				}
				accesses = append(accesses, access)
			}
			// convert []models.permission to []*model.Access
			perm := &update.RobotPermission{
				Kind:      bot.Permissions[0].Kind,
				Namespace: bot.Permissions[0].Namespace,
				Access:    accesses,
			}
			opts.Permissions = []*update.RobotPermission{perm}

			err = updateRobotView(&opts)
			if err != nil {
				log.Errorf("failed to Update robot")
			}
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(
		&all,
		"all-permission",
		"a",
		false,
		"Select all permissions for the robot account",
	)
	flags.StringVarP(&opts.Name, "name", "", "", "name of the robot account")
	flags.StringVarP(&opts.Description, "description", "", "", "description of the robot account")
	flags.Int64VarP(&opts.Duration, "duration", "", 0, "set expiration of robot account in days")

	return cmd
}

func updateRobotView(updateView *update.UpdateView) error {
	if updateView == nil {
		updateView = &update.UpdateView{}
	}

	update.UpdateRobotView(updateView)
	return api.UpdateRobot(updateView)
}
