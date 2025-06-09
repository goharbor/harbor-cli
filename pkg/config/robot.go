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
package config

import (
	"fmt"
	"os"
	"slices"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	"gopkg.in/yaml.v2"
)

type RobotPermissionConfig struct {
	Name        string           `yaml:"name"`
	Description string           `yaml:"description"`
	Duration    int64            `yaml:"duration"`
	Project     string           `yaml:"project"`
	Permissions []PermissionSpec `yaml:"permissions"`
}

type PermissionSpec struct {
	Resource  string   `yaml:"resource,omitempty"`
	Resources []string `yaml:"resources,omitempty"`
	Actions   []string `yaml:"actions"`
}

func LoadRobotConfigFromYAML(filename string) (*create.CreateView, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %v", err)
	}

	var config RobotPermissionConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}

	opts := &create.CreateView{
		Name:        config.Name,
		Description: config.Description,
		Duration:    config.Duration,
		ProjectName: config.Project,
	}

	permissions, err := ProcessPermissions(config.Permissions)
	if err != nil {
		return nil, err
	}

	var accesses []*models.Access
	for _, perm := range permissions {
		access := &models.Access{
			Action:   perm.Action,
			Resource: perm.Resource,
		}
		accesses = append(accesses, access)
	}

	perm := &create.RobotPermission{
		Namespace: config.Project,
		Access:    accesses,
	}
	opts.Permissions = []*create.RobotPermission{perm}

	return opts, nil
}

func ProcessPermissions(specs []PermissionSpec) ([]models.Permission, error) {
	var result []models.Permission

	availablePerms, err := GetAllAvailablePermissions()
	if err != nil {
		return nil, err
	}

	for _, spec := range specs {
		var resources []string

		if spec.Resource != "" {
			resources = []string{spec.Resource}
		} else if len(spec.Resources) > 0 {
			resources = spec.Resources
		} else {
			return nil, fmt.Errorf("permission must specify either 'resource' or 'resources'")
		}

		if containsWildcard(resources) {
			resources = getAllResourceNames(availablePerms)
		}

		for _, resource := range resources {
			if !isValidResource(resource, availablePerms) && resource != "*" {
				fmt.Printf("Warning: Resource '%s' is not valid and will be skipped\n", resource)
				continue
			}

			if containsWildcard(spec.Actions) {
				validActions := getValidActionsForResource(resource, availablePerms)
				for _, action := range validActions {
					result = append(result, models.Permission{
						Resource: resource,
						Action:   action,
					})
				}
			} else {
				for _, action := range spec.Actions {
					if isValidAction(resource, action, availablePerms) {
						result = append(result, models.Permission{
							Resource: resource,
							Action:   action,
						})
					} else {
						fmt.Printf("Warning: Action '%s' is not valid for resource '%s' and will be skipped\n",
							action, resource)
					}
				}
			}
		}
	}

	return result, nil
}

func LoadRobotConfigFromFile(filename string) (*create.CreateView, error) {
	opts, err := LoadRobotConfigFromYAML(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}
	if opts.Name == "" {
		return nil, fmt.Errorf("robot name cannot be empty")
	}
	if opts.Duration == 0 {
		return nil, fmt.Errorf("duration cannot be 0")
	}
	if opts.ProjectName == "" {
		return nil, fmt.Errorf("project name cannot be empty")
	}
	if len(opts.Permissions) == 0 || len(opts.Permissions[0].Access) == 0 {
		return nil, fmt.Errorf("no permissions specified")
	}

	projectExists := false
	projectsResp, err := api.ListAllProjects()
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %v", err)
	}

	for _, proj := range projectsResp.Payload {
		if proj.Name == opts.ProjectName {
			projectExists = true
			break
		}
	}

	if !projectExists {
		return nil, fmt.Errorf("project '%s' does not exist in Harbor", opts.ProjectName)
	}
	return opts, nil
}

func GetAllAvailablePermissions() (map[string][]string, error) {
	permsResp, err := api.GetPermissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %v", err)
	}

	result := make(map[string][]string)
	for _, perm := range permsResp.Payload.Project {
		resource := perm.Resource
		if _, exists := result[resource]; !exists {
			result[resource] = []string{}
		}
		result[resource] = append(result[resource], perm.Action)
	}

	return result, nil
}

func containsWildcard(items []string) bool {
	return slices.Contains(items, "*")
}

func getAllResourceNames(permissions map[string][]string) []string {
	resources := make([]string, 0, len(permissions))
	for resource := range permissions {
		resources = append(resources, resource)
	}
	return resources
}

func isValidResource(resource string, permissions map[string][]string) bool {
	_, exists := permissions[resource]
	return exists
}

func isValidAction(resource, action string, permissions map[string][]string) bool {
	actions, exists := permissions[resource]
	if !exists {
		return false
	}

	return slices.Contains(actions, action)
}

func getValidActionsForResource(resource string, permissions map[string][]string) []string {
	if actions, exists := permissions[resource]; exists {
		return actions
	}
	return []string{}
}
