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
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/sirupsen/logrus"
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

// LoadConfigFromFile reads and validates a robot spec from a YAML or JSON
// file.
func LoadConfigFromFile(filename string) (*robotview.CreateView, error) {
	view, err := loadConfigFromYAMLorJSON(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}
	if err := validateLoadedView(view); err != nil {
		return nil, err
	}
	return view, nil
}

func loadConfigFromYAMLorJSON(filename string) (*robotview.CreateView, error) {
	config, err := parseConfigFile(filename)
	if err != nil {
		return nil, err
	}

	level := "project" // Default to project robot
	if config.Level != "" {
		level = strings.ToLower(config.Level)
	}

	permissions, err := expandPermissionScopes(config.Permissions, level)
	if err != nil {
		return nil, err
	}
	if len(permissions) == 0 {
		return nil, fmt.Errorf("no permissions defined in the configuration")
	}

	view := &robotview.CreateView{
		Name:        config.Name,
		Description: config.Description,
		Duration:    config.Duration,
		Level:       level,
		Permissions: permissions,
	}
	if level == "project" {
		view.ProjectName = permissions[0].Namespace
	}
	return view, nil
}

func parseConfigFile(filename string) (*RobotPermissionConfig, error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		return nil, fmt.Errorf("file must have an extension (.yaml, .yml, or .json)")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	var config RobotPermissionConfig
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %v", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %v", err)
		}
	default:
		return nil, fmt.Errorf("unsupported file type: %s, expected 'yaml' or 'json'", ext[1:])
	}
	return &config, nil
}

// expandPermissionScopes resolves wildcards and validates every scope
// against the permissions Harbor offers. Invalid resources and actions are
// skipped with a warning.
func expandPermissionScopes(scopes []PermissionScope, level string) ([]*robotview.RobotPermission, error) {
	available, err := GetAllAvailablePermissions()
	if err != nil {
		return nil, err
	}

	var result []*robotview.RobotPermission
	for _, scope := range scopes {
		expanded, err := expandScope(scope, level, available)
		if err != nil {
			return nil, err
		}
		if len(expanded.Access) > 0 {
			result = append(result, expanded)
		}
	}

	if level == "project" && len(result) > 1 {
		return nil, fmt.Errorf("project robots can only have one permission scope, found %d", len(result))
	}
	return result, nil
}

func expandScope(scope PermissionScope, level string, available *PermissionMap) (*robotview.RobotPermission, error) {
	kind := strings.ToLower(scope.Kind)
	if level == "project" && kind != "project" {
		return nil, fmt.Errorf("project robots can only have project permission scopes, found %s", kind)
	}

	valid := available.System
	if kind == "project" {
		valid = available.Project
	}

	permission := &robotview.RobotPermission{
		Kind:      kind,
		Namespace: scope.Namespace,
		Access:    []*models.Access{},
	}
	for _, item := range scope.Access {
		access, err := expandAccessItem(item, kind, valid)
		if err != nil {
			return nil, err
		}
		permission.Access = append(permission.Access, access...)
	}
	return permission, nil
}

func expandAccessItem(item AccessItem, kind string, valid map[string][]string) ([]*models.Access, error) {
	resources, err := resolveResources(item, valid)
	if err != nil {
		return nil, err
	}

	var access []*models.Access
	for _, resource := range resources {
		if !isValidResource(resource, valid) && resource != "*" {
			logrus.Warnf("Resource '%s' is not valid for scope '%s' and will be skipped", resource, kind)
			continue
		}
		access = append(access, expandActions(resource, item.Actions, kind, valid)...)
	}
	return access, nil
}

func resolveResources(item AccessItem, valid map[string][]string) ([]string, error) {
	var resources []string
	switch {
	case item.Resource != "":
		resources = []string{item.Resource}
	case len(item.Resources) > 0:
		resources = item.Resources
	default:
		return nil, fmt.Errorf("permission must specify either 'resource' or 'resources'")
	}

	if containsWildcard(resources) {
		return getAllResourceNames(valid), nil
	}
	return resources, nil
}

func expandActions(resource string, actions []string, kind string, valid map[string][]string) []*models.Access {
	if containsWildcard(actions) {
		actions = getValidActionsForResource(resource, valid)
	}

	var access []*models.Access
	for _, action := range actions {
		if !isValidAction(resource, action, valid) {
			logrus.Warnf("Action '%s' is not valid for resource '%s' in scope '%s' and will be skipped",
				action, resource, kind)
			continue
		}
		access = append(access, &models.Access{Resource: resource, Action: action})
	}
	return access
}

func validateLoadedView(view *robotview.CreateView) error {
	if view.Name == "" {
		return fmt.Errorf("robot name cannot be empty")
	}
	if view.Duration == 0 {
		return fmt.Errorf("duration cannot be 0")
	}
	if len(view.Permissions) == 0 {
		return fmt.Errorf("no permissions specified")
	}
	for _, permission := range view.Permissions {
		if len(permission.Access) == 0 {
			return fmt.Errorf("no access defined for permission scope")
		}
	}

	switch view.Level {
	case "project":
		return validateProjectRobotView(view)
	case "system":
		return nil
	default:
		return fmt.Errorf("invalid robot level: %s. Must be 'project' or 'system'", view.Level)
	}
}

func validateProjectRobotView(view *robotview.CreateView) error {
	if containsWildcard([]string{view.Permissions[0].Namespace}) {
		return fmt.Errorf("project robots cannot have wildcard namespace in permissions")
	}

	projectName := view.ProjectName
	if projectName == "" {
		projectName = view.Permissions[0].Namespace
	}
	if projectName == "" {
		return fmt.Errorf("project name is required for project robots")
	}

	if _, err := api.GetProject(projectName, false); err != nil {
		return fmt.Errorf("failed to get project: %v, does it exist?", err)
	}

	view.ProjectName = projectName
	return nil
}

// GetAllAvailablePermissions returns permissions organized by scope. The
// error is returned unwrapped so callers can still parse the Harbor
// response for its status and message.
func GetAllAvailablePermissions() (*PermissionMap, error) {
	response, err := api.GetPermissions()
	if err != nil {
		return nil, err
	}

	result := &PermissionMap{
		Project: make(map[string][]string),
		System:  make(map[string][]string),
	}
	for _, permission := range response.Payload.Project {
		result.Project[permission.Resource] = append(result.Project[permission.Resource], permission.Action)
	}
	for _, permission := range response.Payload.System {
		result.System[permission.Resource] = append(result.System[permission.Resource], permission.Action)
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
	return slices.Contains(permissions[resource], action)
}

func getValidActionsForResource(resource string, permissions map[string][]string) []string {
	return permissions[resource]
}
