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
	config "github.com/goharbor/harbor-cli/pkg/config/robot"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	"github.com/goharbor/harbor-cli/pkg/views/robot/update"
	"github.com/sirupsen/logrus"
)

func loadRobotConfigFromFile(configFile string, permissions *[]models.Permission, projectPermissionsMap map[string][]models.Permission, isUpdate bool, createOpts *create.CreateView, updateOpts *update.UpdateView) error {
	fmt.Println("Loading configuration from: ", configFile)

	loadedOpts, err := config.LoadRobotConfigFromFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to load robot config from file: %v", err)
	}

	logrus.Info("Successfully loaded robot configuration")

	// Apply configuration based on operation type
	if isUpdate && updateOpts != nil {
		// Update mode: Only update specific fields selectively
		if loadedOpts.Description != "" {
			updateOpts.Description = loadedOpts.Description
		}
		if loadedOpts.Duration != 0 {
			updateOpts.Duration = loadedOpts.Duration
		}
	} else if !isUpdate && createOpts != nil {
		// Create mode: Full assignment
		*createOpts = *loadedOpts
	}

	var systemPermFound bool
	for _, perm := range loadedOpts.Permissions {
		if perm.Kind == "system" && perm.Namespace == "/" {
			systemPermFound = true

			if isUpdate {
				// Append for updates
				for _, access := range perm.Access {
					*permissions = append(*permissions, models.Permission{
						Resource: access.Resource,
						Action:   access.Action,
					})
				}
			} else {
				// Replace for creates
				*permissions = make([]models.Permission, len(perm.Access))
				for i, access := range perm.Access {
					(*permissions)[i] = models.Permission{
						Resource: access.Resource,
						Action:   access.Action,
					}
				}
			}
		} else if perm.Kind == "project" {
			var projectPerms []models.Permission
			for _, access := range perm.Access {
				projectPerms = append(projectPerms, models.Permission{
					Resource: access.Resource,
					Action:   access.Action,
				})
			}
			projectPermissionsMap[perm.Namespace] = projectPerms
		}
	}

	if !systemPermFound && loadedOpts.Level == "system" {
		return fmt.Errorf("robot configuration must include system-level permissions")
	}

	logrus.Infof("Loaded robot config with %d system permissions and %d project-specific permissions",
		len(*permissions), len(projectPermissionsMap))

	return nil
}

func LoadFromConfigFileForCreate(opts *create.CreateView, configFile string, permissions *[]models.Permission, projectPermissionsMap map[string][]models.Permission) error {
	return loadRobotConfigFromFile(configFile, permissions, projectPermissionsMap, false, opts, nil)
}

func LoadFromConfigFileForUpdate(opts *update.UpdateView, configFile string, permissions *[]models.Permission, projectPermissionsMap map[string][]models.Permission) error {
	return loadRobotConfigFromFile(configFile, permissions, projectPermissionsMap, true, nil, opts)
}
