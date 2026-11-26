// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package robot

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// UpdateOptions parameterizes RunUpdate for both robot levels.
type UpdateOptions struct {
	Level      Level
	RobotID    int64
	All        bool
	ConfigFile string
}

// RunUpdate loads a robot, applies changes from a config file or
// interactive input, and saves it.
func RunUpdate(opts UpdateOptions) error {
	robot, err := api.GetRobot(opts.RobotID)
	if err != nil {
		return WrapError("get robot", err)
	}

	view := viewFromRobot(robot.Payload)
	permissions := NewPermissionSet()
	permissions.Split(robot.Payload.Permissions)
	logrus.Debugf("Loaded robot with %d system permissions and %d project-specific permissions",
		len(permissions.System), len(permissions.Project))

	if opts.ConfigFile != "" {
		if permissions, err = ApplyConfigToUpdate(&view, opts.Level, opts.ConfigFile); err != nil {
			return err
		}
	} else if err := promptUpdate(&view, opts.Level, opts.All, permissions); err != nil {
		return err
	}

	view.Permissions = permissions.Merge()
	if len(view.Permissions) == 0 {
		return errNoPermissions("update")
	}

	if err := api.UpdateRobot(&view); err != nil {
		return WrapError("update robot", err)
	}
	return showUpdatedRobot(&view)
}

// viewFromRobot copies a robot into an UpdateView, guarding the optional
// fields.
func viewFromRobot(robot *models.Robot) robotview.UpdateView {
	view := robotview.UpdateView{
		CreationTime: robot.CreationTime,
		Description:  robot.Description,
		Disable:      robot.Disable,
		Editable:     robot.Editable,
		ID:           robot.ID,
		Level:        robot.Level,
		Name:         robot.Name,
		Secret:       robot.Secret,
	}
	if robot.Duration != nil {
		view.Duration = *robot.Duration
	}
	return view
}

func showUpdatedRobot(view *robotview.UpdateView) error {
	fmt.Printf("Successfully updated robot account '%s' (ID: %d)\n", view.Name, view.ID)

	if viper.GetString("output-format") != "" {
		details, err := api.GetRobot(view.ID)
		if err != nil {
			return fmt.Errorf("failed to get robot details: %v", utils.ParseHarborErrorMsg(err))
		}
		utils.SavePayloadJSON(view.Name, details.Payload)
	}
	return nil
}
