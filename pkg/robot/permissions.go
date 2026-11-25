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
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/sirupsen/logrus"
)

// PermissionSet separates a robot's permissions into the system scope and
// one scope per project, the shape the interactive flows work with.
type PermissionSet struct {
	System  []models.Permission
	Project map[string][]models.Permission
}

func NewPermissionSet() *PermissionSet {
	return &PermissionSet{Project: make(map[string][]models.Permission)}
}

// Split buckets a robot's permission scopes into the set.
func (ps *PermissionSet) Split(scopes []*robotview.RobotPermission) {
	for _, scope := range scopes {
		if scope.Kind == "system" && scope.Namespace == "/" {
			ps.System = append(ps.System, toPermissions(scope.Access)...)
		} else if scope.Kind == "project" {
			ps.Project[scope.Namespace] = toPermissions(scope.Access)
		}
	}
}

// Merge rebuilds the permission scopes the Harbor API expects. Empty scopes
// are dropped because Harbor rejects them.
func (ps *PermissionSet) Merge() []*robotview.RobotPermission {
	var merged []*robotview.RobotPermission

	for projectName, permissions := range ps.Project {
		if access := toAccess(permissions); len(access) > 0 {
			merged = append(merged, &robotview.RobotPermission{
				Namespace: projectName,
				Access:    access,
				Kind:      "project",
			})
		}
	}

	if access := toAccess(ps.System); len(access) > 0 {
		merged = append(merged, &robotview.RobotPermission{
			Namespace: "/",
			Access:    access,
			Kind:      "system",
		})
	}

	return merged
}

// IsEmpty reports whether no scope holds any permission.
func (ps *PermissionSet) IsEmpty() bool {
	if len(ps.System) > 0 {
		return false
	}
	for _, permissions := range ps.Project {
		if len(permissions) > 0 {
			return false
		}
	}
	return true
}

// ClearProjects drops all project scopes.
func (ps *PermissionSet) ClearProjects() {
	ps.Project = make(map[string][]models.Permission)
}

// ProjectNames lists the projects the set has scopes for.
func (ps *PermissionSet) ProjectNames() []string {
	names := make([]string, 0, len(ps.Project))
	for name := range ps.Project {
		names = append(names, name)
	}
	return names
}

func toAccess(permissions []models.Permission) []*models.Access {
	var access []*models.Access
	for _, permission := range permissions {
		access = append(access, &models.Access{
			Resource: permission.Resource,
			Action:   permission.Action,
		})
	}
	return access
}

func toPermissions(access []*models.Access) []models.Permission {
	var permissions []models.Permission
	for _, a := range access {
		permissions = append(permissions, models.Permission{
			Resource: a.Resource,
			Action:   a.Action,
		})
	}
	return permissions
}

// AllPermissions fetches every permission Harbor offers for the given kind
// ("system" or "project").
func AllPermissions(kind string) ([]models.Permission, error) {
	response, err := api.GetPermissions()
	if err != nil {
		return nil, WrapError("get permissions", err)
	}

	available := response.Payload.Project
	if kind == "system" {
		available = response.Payload.System
	}

	permissions := make([]models.Permission, 0, len(available))
	for _, permission := range available {
		permissions = append(permissions, *permission)
	}
	return permissions, nil
}

// ValidateProjectPermissions filters out permissions Harbor does not offer
// at project level and warns about what was dropped.
func ValidateProjectPermissions(permissions []models.Permission) ([]models.Permission, error) {
	available, err := AllPermissions("project")
	if err != nil {
		return nil, err
	}

	valid := make(map[string]bool, len(available))
	for _, permission := range available {
		valid[permission.Resource+":"+permission.Action] = true
	}

	var kept []models.Permission
	var dropped []string
	for _, permission := range permissions {
		key := permission.Resource + ":" + permission.Action
		if valid[key] {
			kept = append(kept, permission)
		} else {
			dropped = append(dropped, key)
		}
	}

	if len(dropped) > 0 {
		logrus.Warnf("Removed %d invalid project permissions: %v", len(dropped), dropped)
	}
	return kept, nil
}

// errNoPermissions is the uniform complaint when a robot would end up
// without a single permission.
func errNoPermissions(action string) error {
	return fmt.Errorf("failed to %s robot: no permissions selected, robot account needs at least one permission", action)
}
