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
	baseprompt "github.com/goharbor/harbor-cli/pkg/prompt"
	robotprompt "github.com/goharbor/harbor-cli/pkg/prompt/robot"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
)

// GetSystemPermissions consolidates system-permission selection for both create and update flows.
// Behavior:
//   - When isUpdate is true and skipConfirm is false, asks the user whether to update system permissions.
//     If the user chooses No, it keeps existing permissions and returns nil.
//   - When all is true, selects all system permissions from the API and replaces the current slice.
//   - Otherwise, prompts the user to select system permissions; requires at least one.
func GetSystemPermissions(isUpdate bool, skipConfirm bool, all bool, permissions *[]models.Permission) error {
	// Optional confirmation in update flows
	if isUpdate && !skipConfirm {
		updateSystem, err := robotprompt.AskUpdateSystemPerms()
		if err != nil {
			return fmt.Errorf("error asking about system permission updates: %v", err)
		}
		if !updateSystem {
			logrus.Info("Keeping existing system permissions")
			return nil
		}
	}

	if all {
		perms, _ := api.GetPermissions()
		// Replace current permissions with all system permissions
		*permissions = (*permissions)[:0]
		for _, perm := range perms.Payload.System {
			*permissions = append(*permissions, *perm)
		}
		return nil
	}

	// Prompt for explicit selection
	newPermissions := baseprompt.GetRobotPermissionsFromUser("system")
	if len(newPermissions) == 0 {
		ctx := "create"
		if isUpdate {
			ctx = "update"
		}
		return fmt.Errorf("failed to %s robot: %v",
			ctx, utils.ParseHarborErrorMsg(fmt.Errorf("no permissions selected, robot account needs at least one permission")))
	}
	*permissions = newPermissions
	return nil
}

func PermissionsToAccess(permissions []models.Permission) []*models.Access {
	var accessesSystem []*models.Access
	for _, perm := range permissions {
		accessesSystem = append(accessesSystem, &models.Access{
			Resource: perm.Resource,
			Action:   perm.Action,
		})
	}
	return accessesSystem
}
