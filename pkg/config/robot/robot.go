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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	rmodel "github.com/goharbor/harbor-cli/pkg/models/robot"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	"gopkg.in/yaml.v2"
)

// PermissionMap holds permissions separated by scope
type PermissionMap struct {
	Project map[string][]string
	System  map[string][]string
}

// RobotPermissionConfig represents the robot account configuration from file
type RobotPermissionConfig struct {
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description" json:"description"`
	Duration    int64             `yaml:"duration" json:"duration"`
	Level       string            `yaml:"level,omitempty" json:"level,omitempty"` // "project" or "system"
	Permissions []PermissionScope `yaml:"permissions" json:"permissions"`
}

// PermissionScope represents a permission scope with access items, kind and namespace
type PermissionScope struct {
	Access    []AccessItem `yaml:"access" json:"access"`
	Kind      string       `yaml:"kind" json:"kind"`           // "project" or "system"
	Namespace string       `yaml:"namespace" json:"namespace"` // Project name or "/"
}

// AccessItem represents a resource permission definition
type AccessItem struct {
	Resource  string   `yaml:"resource,omitempty" json:"resource,omitempty"`
	Resources []string `yaml:"resources,omitempty" json:"resources,omitempty"`
	Actions   []string `yaml:"actions" json:"actions"`
}

type RobotSecret struct {
	Name         string `json:"name"`
	ExpiresAt    int64  `json:"expires_at"`
	CreationTime string `json:"creation_time"`
	Secret       string `json:"secret"`
}

func LoadRobotConfigFromYAMLorJSON(filename string, fileType string) (*create.CreateView, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var config RobotPermissionConfig
	if fileType == "yaml" {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %v", err)
		}
	} else if fileType == "json" {
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported file type: %s, expected 'yaml' or 'json'", fileType)
	}

	// Determine the robot level and ensure it's lowercase
	robotLevel := "project" // Default to project robot
	if config.Level != "" {
		robotLevel = strings.ToLower(config.Level)
	}

	// Create the base view object
	opts := &create.CreateView{
		Name:        config.Name,
		Description: config.Description,
		Duration:    config.Duration,
		Level:       robotLevel, // Keep lowercase
	}

	// Process permissions
	robotPermissions, err := processPermissionScopes(config.Permissions, opts.Level)
	if err != nil {
		return nil, err
	}

	if len(robotPermissions) == 0 {
		return nil, fmt.Errorf("no permissions defined in the configuration")
	}

	opts.Permissions = robotPermissions

	// If this is a project robot, set the ProjectName from the first permission's namespace
	if opts.Level == "project" && len(opts.Permissions) > 0 {
		opts.ProjectName = opts.Permissions[0].Namespace
	}

	return opts, nil
}

// Process permission scopes
func processPermissionScopes(scopes []PermissionScope, robotLevel string) ([]*rmodel.RobotPermission, error) {
	var result []*rmodel.RobotPermission

	// Get available permissions for validation
	availablePerms, err := GetAllAvailablePermissions()
	if err != nil {
		return nil, err
	}

	// Check if we're creating a project robot
	isProjectRobot := robotLevel == "project"

	// Process each permission scope
	for _, scope := range scopes {
		// Make sure kind is lowercase
		scopeKind := strings.ToLower(scope.Kind)

		// For project robots, enforce only project permissions
		if isProjectRobot && scopeKind != "project" {
			return nil, fmt.Errorf("project robots can only have project permission scopes, found %s", scopeKind)
		}

		robotPerm := &rmodel.RobotPermission{
			Kind:      scopeKind,
			Namespace: scope.Namespace,
			Access:    []*models.Access{},
		}

		// Process access items in this scope
		for _, accessItem := range scope.Access {
			// Get the list of resources
			var resources []string
			if accessItem.Resource != "" {
				resources = []string{accessItem.Resource}
			} else if len(accessItem.Resources) > 0 {
				resources = accessItem.Resources
			} else {
				return nil, fmt.Errorf("permission must specify either 'resource' or 'resources'")
			}

			// Handle wildcard for resources
			if containsWildcard(resources) {
				// Get appropriate resources based on scope kind
				if scopeKind == "project" {
					resources = getAllResourceNames(availablePerms.Project)
				} else {
					resources = getAllResourceNames(availablePerms.System)
				}
			}

			// Process each resource
			for _, resource := range resources {
				// Determine which permission map to use based on scope kind
				var validActions []string
				var isValid bool

				if scopeKind == "project" {
					isValid = isValidResource(resource, availablePerms.Project)
					validActions = getValidActionsForResource(resource, availablePerms.Project)
				} else {
					isValid = isValidResource(resource, availablePerms.System)
					validActions = getValidActionsForResource(resource, availablePerms.System)
				}

				if !isValid && resource != "*" {
					fmt.Printf("Warning: Resource '%s' is not valid for scope '%s' and will be skipped\n",
						resource, scopeKind)
					continue
				}

				// Handle wildcard for actions
				if containsWildcard(accessItem.Actions) {
					for _, action := range validActions {
						robotPerm.Access = append(robotPerm.Access, &models.Access{
							Resource: resource,
							Action:   action,
						})
					}
				} else {
					// Process specific actions
					for _, action := range accessItem.Actions {
						var actionValid bool

						if scopeKind == "project" {
							actionValid = isValidAction(resource, action, availablePerms.Project)
						} else {
							actionValid = isValidAction(resource, action, availablePerms.System)
						}

						if actionValid {
							robotPerm.Access = append(robotPerm.Access, &models.Access{
								Resource: resource,
								Action:   action,
							})
						} else {
							fmt.Printf("Warning: Action '%s' is not valid for resource '%s' in scope '%s' and will be skipped\n",
								action, resource, scopeKind)
						}
					}
				}
			}
		}

		// Add the permission scope if it has any access items
		if len(robotPerm.Access) > 0 {
			result = append(result, robotPerm)
		}
	}

	// For project robots, ensure there's only one permission scope
	if isProjectRobot && len(result) > 1 {
		return nil, fmt.Errorf("project robots can only have one permission scope, found %d", len(result))
	}

	return result, nil
}

func LoadRobotConfigFromFile(filename string) (*create.CreateView, error) {
	var opts *create.CreateView
	var err error

	ext := filepath.Ext(filename)
	if ext == "" {
		return nil, fmt.Errorf("file must have an extension (.yaml, .yml, or .json)")
	}

	fileType := ext[1:] // Remove the dot
	if fileType == "yml" {
		fileType = "yaml"
	}

	opts, err = LoadRobotConfigFromYAMLorJSON(filename, fileType)

	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	// Basic validation
	if opts.Name == "" {
		return nil, fmt.Errorf("robot name cannot be empty")
	}
	if opts.Duration == 0 {
		return nil, fmt.Errorf("duration cannot be 0")
	}

	// Level-specific validation
	if opts.Level == "project" {
		// Project robot requires a project
		projectName := ""
		if opts.ProjectName != "" {
			projectName = opts.ProjectName
		} else if len(opts.Permissions) > 0 && opts.Permissions[0].Namespace != "" {
			projectName = opts.Permissions[0].Namespace
		}

		if projectName == "" {
			return nil, fmt.Errorf("project name is required for project robots")
		}

		// Verify project exists
		projectExists := false
		projectsResp, err := api.ListAllProjects()
		if err != nil {
			return nil, fmt.Errorf("failed to list projects: %v", err)
		}

		for _, proj := range projectsResp.Payload {
			if proj.Name == projectName {
				projectExists = true
				break
			}
		}

		if !projectExists {
			return nil, fmt.Errorf("project '%s' does not exist in Harbor", projectName)
		}

		// Set the project name consistently
		opts.ProjectName = projectName
	} else if opts.Level == "system" {
		// System robot validation
		if len(opts.Permissions) == 0 {
			return nil, fmt.Errorf("system robot must have at least one permission scope")
		}
	} else {
		return nil, fmt.Errorf("invalid robot level: %s. Must be 'project' or 'system'", opts.Level)
	}

	// Validate permissions
	if len(opts.Permissions) == 0 {
		return nil, fmt.Errorf("no permissions specified")
	}

	for _, perm := range opts.Permissions {
		if len(perm.Access) == 0 {
			return nil, fmt.Errorf("no access defined for permission scope")
		}
	}

	return opts, nil
}

// GetAllAvailablePermissions returns permissions organized by scope
func GetAllAvailablePermissions() (*PermissionMap, error) {
	permsResp, err := api.GetPermissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %v", err)
	}

	result := &PermissionMap{
		Project: make(map[string][]string),
		System:  make(map[string][]string),
	}

	// Add project permissions
	for _, perm := range permsResp.Payload.Project {
		resource := perm.Resource
		if _, exists := result.Project[resource]; !exists {
			result.Project[resource] = []string{}
		}
		result.Project[resource] = append(result.Project[resource], perm.Action)
	}

	// Add system permissions if available
	if permsResp.Payload.System != nil {
		for _, perm := range permsResp.Payload.System {
			resource := perm.Resource
			if _, exists := result.System[resource]; !exists {
				result.System[resource] = []string{}
			}
			result.System[resource] = append(result.System[resource], perm.Action)
		}
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
