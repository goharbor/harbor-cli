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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	rmodel "github.com/goharbor/harbor-cli/pkg/models/robot"
)

func buildMergedPermissions(projectPermissionsMap map[string][]models.Permission, accessesSystem []*models.Access) []*rmodel.RobotPermission {
	var mergedPermissions []*rmodel.RobotPermission

	// Add project permissions
	for projectName, projectPermissions := range projectPermissionsMap {
		var accessesProject []*models.Access
		for _, perm := range projectPermissions {
			accessesProject = append(accessesProject, &models.Access{
				Resource: perm.Resource,
				Action:   perm.Action,
			})
		}
		mergedPermissions = append(mergedPermissions, &rmodel.RobotPermission{
			Namespace: projectName,
			Access:    accessesProject,
			Kind:      "project",
		})
	}

	// Add system permissions
	mergedPermissions = append(mergedPermissions, &rmodel.RobotPermission{
		Namespace: "/",
		Access:    accessesSystem,
		Kind:      "system",
	})

	return mergedPermissions
}
