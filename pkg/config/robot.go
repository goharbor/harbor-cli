// // Copyright Project Harbor Authors
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// // http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.
// package config

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"slices"

// 	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
// 	"github.com/goharbor/harbor-cli/pkg/api"
// 	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
// 	"gopkg.in/yaml.v2"
// )

// type RobotPermissionConfig struct {
// 	Name        string           `yaml:"name" json:"name"`
// 	Description string           `yaml:"description" json:"description"`
// 	Duration    int64            `yaml:"duration" json:"duration"`
// 	Project     string           `yaml:"project" json:"project"`
// 	Permissions []PermissionSpec `yaml:"permissions" json:"permissions"`
// }

// type PermissionSpec struct {
// 	Resource  string   `yaml:"resource,omitempty" json:"resource,omitempty"`
// 	Resources []string `yaml:"resources,omitempty" json:"resources,omitempty"`
// 	Actions   []string `yaml:"actions" json:"actions"`
// }

// type RobotSecret struct {
// 	Name         string `json:"name"`
// 	ExpiresAt    int64  `json:"expires_at"`
// 	CreationTime string `json:"creation_time"`
// 	Secret       string `json:"secret"`
// }

// func LoadRobotConfigFromYAMLorJSON(filename string, fileType string) (*create.CreateView, error) {
// 	data, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read YAML file: %v", err)
// 	}
// 	var config RobotPermissionConfig
// 	if fileType == "yaml" {
// 		if err := yaml.Unmarshal(data, &config); err != nil {
// 			return nil, fmt.Errorf("failed to parse YAML: %v", err)
// 		}
// 	} else if fileType == "json" {
// 		if err := json.Unmarshal(data, &config); err != nil {
// 			return nil, fmt.Errorf("failed to parse JSON: %v", err)
// 		}
// 	} else {
// 		return nil, fmt.Errorf("unsupported file type: %s, expected 'yaml' or 'json'", fileType)
// 	}

// 	opts := &create.CreateView{
// 		Name:        config.Name,
// 		Description: config.Description,
// 		Duration:    config.Duration,
// 		ProjectName: config.Project,
// 	}

// 	permissions, err := ProcessPermissions(config.Permissions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var accesses []*models.Access
// 	for _, perm := range permissions {
// 		access := &models.Access{
// 			Action:   perm.Action,
// 			Resource: perm.Resource,
// 		}
// 		accesses = append(accesses, access)
// 	}

// 	perm := &create.RobotPermission{
// 		Namespace: config.Project,
// 		Access:    accesses,
// 	}
// 	opts.Permissions = []*create.RobotPermission{perm}

// 	return opts, nil
// }

// func ProcessPermissions(specs []PermissionSpec) ([]models.Permission, error) {
// 	var result []models.Permission

// 	availablePerms, err := GetAllAvailablePermissions()
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, spec := range specs {
// 		var resources []string

// 		if spec.Resource != "" {
// 			resources = []string{spec.Resource}
// 		} else if len(spec.Resources) > 0 {
// 			resources = spec.Resources
// 		} else {
// 			return nil, fmt.Errorf("permission must specify either 'resource' or 'resources'")
// 		}

// 		if containsWildcard(resources) {
// 			resources = getAllResourceNames(availablePerms)
// 		}

// 		for _, resource := range resources {
// 			if !isValidResource(resource, availablePerms) && resource != "*" {
// 				fmt.Printf("Warning: Resource '%s' is not valid and will be skipped\n", resource)
// 				continue
// 			}

// 			if containsWildcard(spec.Actions) {
// 				validActions := getValidActionsForResource(resource, availablePerms)
// 				for _, action := range validActions {
// 					result = append(result, models.Permission{
// 						Resource: resource,
// 						Action:   action,
// 					})
// 				}
// 			} else {
// 				for _, action := range spec.Actions {
// 					if isValidAction(resource, action, availablePerms) {
// 						result = append(result, models.Permission{
// 							Resource: resource,
// 							Action:   action,
// 						})
// 					} else {
// 						fmt.Printf("Warning: Action '%s' is not valid for resource '%s' and will be skipped\n",
// 							action, resource)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return result, nil
// }

// func LoadRobotConfigFromFile(filename string) (*create.CreateView, error) {
// 	var opts *create.CreateView
// 	var err error
// 	opts, err = LoadRobotConfigFromYAMLorJSON(filename, filepath.Ext(filename)[1:])

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load configuration: %v", err)
// 	}
// 	if opts.Name == "" {
// 		return nil, fmt.Errorf("robot name cannot be empty")
// 	}
// 	if opts.Duration == 0 {
// 		return nil, fmt.Errorf("duration cannot be 0")
// 	}
// 	if opts.ProjectName == "" {
// 		return nil, fmt.Errorf("project name cannot be empty")
// 	}
// 	if len(opts.Permissions) == 0 || len(opts.Permissions[0].Access) == 0 {
// 		return nil, fmt.Errorf("no permissions specified")
// 	}

// 	projectExists := false
// 	projectsResp, err := api.ListAllProjects()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list projects: %v", err)
// 	}

// 	for _, proj := range projectsResp.Payload {
// 		if proj.Name == opts.ProjectName {
// 			projectExists = true
// 			break
// 		}
// 	}

// 	if !projectExists {
// 		return nil, fmt.Errorf("project '%s' does not exist in Harbor", opts.ProjectName)
// 	}
// 	return opts, nil
// }

// func GetAllAvailablePermissions() (map[string][]string, error) {
// 	permsResp, err := api.GetPermissions()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get permissions: %v", err)
// 	}

// 	result := make(map[string][]string)
// 	for _, perm := range permsResp.Payload.Project {
// 		resource := perm.Resource
// 		if _, exists := result[resource]; !exists {
// 			result[resource] = []string{}
// 		}
// 		result[resource] = append(result[resource], perm.Action)
// 	}

// 	return result, nil
// }

// func containsWildcard(items []string) bool {
// 	return slices.Contains(items, "*")
// }

// func getAllResourceNames(permissions map[string][]string) []string {
// 	resources := make([]string, 0, len(permissions))
// 	for resource := range permissions {
// 		resources = append(resources, resource)
// 	}
// 	return resources
// }

// func isValidResource(resource string, permissions map[string][]string) bool {
// 	_, exists := permissions[resource]
// 	return exists
// }

// func isValidAction(resource, action string, permissions map[string][]string) bool {
// 	actions, exists := permissions[resource]
// 	if !exists {
// 		return false
// 	}

// 	return slices.Contains(actions, action)
// }

// func getValidActionsForResource(resource string, permissions map[string][]string) []string {
// 	if actions, exists := permissions[resource]; exists {
// 		return actions
// 	}
// 	return []string{}
// }

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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// RobotPermissionConfig represents the robot account configuration from file
type RobotPermissionConfig struct {
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description" json:"description"`
	Duration    int64             `yaml:"duration" json:"duration"`
	Project     string            `yaml:"project,omitempty" json:"project,omitempty"` // Legacy field for backward compatibility
	Kind        string            `yaml:"kind,omitempty" json:"kind,omitempty"`       // "project" or "system"
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

// Legacy PermissionSpec for backward compatibility
type PermissionSpec struct {
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

	// Determine the robot kind
	robotKind := "project" // Default to project robot
	if config.Kind != "" {
		robotKind = config.Kind
	}

	// Create the base view object
	opts := &create.CreateView{
		Name:        config.Name,
		Description: config.Description,
		Duration:    config.Duration,
		Kind:        robotKind,
	}

	// For backward compatibility, set ProjectName if present in config
	if config.Project != "" {
		opts.ProjectName = config.Project
	}

	// Process permissions based on the new format
	robotPermissions, err := processPermissionScopes(config.Permissions)
	if err != nil {
		return nil, err
	}

	// If no permissions defined in new format but we have legacy format,
	// try to process them as backward compatibility
	if len(robotPermissions) == 0 && len(config.Permissions) == 0 {
		// This assumes you're passing legacy PermissionSpec items directly to the function
		// You may need to adjust this depending on how the legacy format is actually provided
		return nil, fmt.Errorf("no permissions defined in the configuration")
	}

	opts.Permissions = robotPermissions

	// For backward compatibility with legacy format
	// (if no permission scopes but Project is defined)
	if len(opts.Permissions) == 0 && opts.ProjectName != "" {
		log.Warn("Using legacy format for robot configuration")
		// Create a default permission with the project
		opts.Permissions = []*create.RobotPermission{
			{
				Kind:      "project",
				Namespace: opts.ProjectName,
				Access:    []*models.Access{},
			},
		}
	}

	return opts, nil
}

// Process permission scopes from the new format
func processPermissionScopes(scopes []PermissionScope) ([]*create.RobotPermission, error) {
	var result []*create.RobotPermission

	// Get available permissions for validation
	availablePerms, err := GetAllAvailablePermissions()
	if err != nil {
		return nil, err
	}

	// Process each permission scope
	for _, scope := range scopes {
		robotPerm := &create.RobotPermission{
			Kind:      scope.Kind,
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
				resources = getAllResourceNames(availablePerms)
			}

			// Process each resource
			for _, resource := range resources {
				if !isValidResource(resource, availablePerms) && resource != "*" {
					fmt.Printf("Warning: Resource '%s' is not valid and will be skipped\n", resource)
					continue
				}

				// Handle wildcard for actions
				if containsWildcard(accessItem.Actions) {
					validActions := getValidActionsForResource(resource, availablePerms)
					for _, action := range validActions {
						robotPerm.Access = append(robotPerm.Access, &models.Access{
							Resource: resource,
							Action:   action,
						})
					}
				} else {
					// Process specific actions
					for _, action := range accessItem.Actions {
						if isValidAction(resource, action, availablePerms) {
							robotPerm.Access = append(robotPerm.Access, &models.Access{
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

		// Add the permission scope if it has any access items
		if len(robotPerm.Access) > 0 {
			result = append(result, robotPerm)
		}
	}

	return result, nil
}

// ProcessPermissions processes the legacy permission format
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

	// Kind-specific validation
	if opts.Kind == "project" {
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

		// For project robots, ensure there's exactly one permission scope
		if len(opts.Permissions) > 1 {
			return nil, fmt.Errorf("project robots can only have one permission scope, found %d", len(opts.Permissions))
		}

		// Also validate that the single permission scope is of kind "project"
		if len(opts.Permissions) == 1 && opts.Permissions[0].Kind != "project" {
			return nil, fmt.Errorf("project robots must have a permission scope of kind 'project', found '%s'",
				opts.Permissions[0].Kind)
		}
	} else if opts.Kind == "system" {
		// System robot validation
		if len(opts.Permissions) == 0 {
			return nil, fmt.Errorf("system robot must have at least one permission scope")
		}
	} else {
		return nil, fmt.Errorf("invalid robot kind: %s. Must be 'project' or 'system'", opts.Kind)
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

func GetAllAvailablePermissions() (map[string][]string, error) {
	permsResp, err := api.GetPermissions()
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %v", err)
	}

	result := make(map[string][]string)

	// Add project permissions
	for _, perm := range permsResp.Payload.Project {
		resource := perm.Resource
		if _, exists := result[resource]; !exists {
			result[resource] = []string{}
		}
		result[resource] = append(result[resource], perm.Action)
	}

	// Add system permissions if available
	if permsResp.Payload.System != nil {
		for _, perm := range permsResp.Payload.System {
			resource := perm.Resource
			if _, exists := result[resource]; !exists {
				result[resource] = []string{}
			}
			result[resource] = append(result[resource], perm.Action)
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
