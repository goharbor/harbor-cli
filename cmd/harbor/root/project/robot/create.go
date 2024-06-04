package robot

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateRobotCommand() *cobra.Command {
	var (
		opts        create.CreateView
		projectName string
		all         bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create robot",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			kind := "project"
			opts.Level = kind

			if opts.ProjectName == "" {
				opts.ProjectName = prompt.GetProjectNameFromUser()
				if opts.ProjectName == "" {
					os.Exit(1)
				}
			}

			// to-do handle permission as json submission
			perm := &create.RobotPermission{
				Kind:      kind,
				Namespace: projectName,
			}
			opts.Permissions = []*create.RobotPermission{perm}

			if len(args) == 0 {
				if opts.Name == "" || opts.Duration == 0 {
					create.CreateRobotView(&opts)
				}
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
				perm := &create.RobotPermission{
					Kind:      kind,
					Namespace: projectName,
					Access:    accesses,
				}
				opts.Permissions = []*create.RobotPermission{perm}
			}
			response, err := api.CreateRobot(opts)
			if err != nil {
				log.Errorf("failed to create robot: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				name := response.Payload.Name
				res, _ := api.GetRobot(response.Payload.ID)
				utils.SavePayloadJSON(name, res.Payload)
				return
			}

			name, secret := response.Payload.Name, response.Payload.Secret
			create.CreateRobotSecretView(name, secret)
			err = clipboard.WriteAll(response.Payload.Secret)
			fmt.Println("secret copied to clipboard.")
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
	flags.StringVarP(&opts.ProjectName, "project", "", "", "set project name")
	flags.StringVarP(&opts.Name, "name", "", "", "name of the robot account")
	flags.StringVarP(&opts.Description, "description", "", "", "description of the robot account")
	flags.Int64VarP(&opts.Duration, "duration", "", 0, "set expiration of robot account in days")

	return cmd
}
