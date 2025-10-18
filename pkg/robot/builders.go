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
	"slices"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	rmodel "github.com/goharbor/harbor-cli/pkg/models/robot"
)

func BuildPermissions(projectPermissionsMap map[string][]models.Permission, accessesSystem ...[]*models.Access) []*rmodel.RobotPermission {
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
	// Add system permissions if any
	if len(accessesSystem) > 0 {
		mergedPermissions = append(mergedPermissions, &rmodel.RobotPermission{
			Namespace: "/",
			Access:    accessesSystem[0],
			Kind:      "system",
		})
	}
	return mergedPermissions
}

// RobotBuilder accumulates fields and permissions for building robot requests.
type RobotBuilder struct {
	name        string
	description string
	duration    int64
	disable     bool

	// Permissions
	projects map[string][]models.Permission
	system   []models.Permission
}

// NewRobotBuilder creates a fresh builder.
func NewRobotBuilder() *RobotBuilder {
	return &RobotBuilder{
		duration: -1, // default aligns with existing CLI UX for "no expiration"
		projects: make(map[string][]models.Permission),
	}
}

func (b *RobotBuilder) WithName(name string) *RobotBuilder {
	b.name = name
	return b
}

func (b *RobotBuilder) WithDescription(desc string) *RobotBuilder {
	b.description = desc
	return b
}

func (b *RobotBuilder) WithDuration(days int64) *RobotBuilder {
	b.duration = days
	return b
}

func (b *RobotBuilder) WithDisable(disable bool) *RobotBuilder {
	b.disable = disable
	return b
}

// Permission setters
func (b *RobotBuilder) AddProjectPermissions(project string, perms ...models.Permission) *RobotBuilder {
	if _, ok := b.projects[project]; !ok {
		b.projects[project] = make([]models.Permission, 0, len(perms))
	}
	b.projects[project] = append(b.projects[project], perms...)
	return b
}

func (b *RobotBuilder) SetProjectPermissions(project string, perms []models.Permission) *RobotBuilder {
	// overwrite
	cp := slices.Clone(perms)
	b.projects[project] = cp
	return b
}

func (b *RobotBuilder) SetSystemPermissions(perms []models.Permission) *RobotBuilder {
	b.system = slices.Clone(perms)
	return b
}

func (b *RobotBuilder) MergePermissions() []*rmodel.RobotPermission {
	var mergedPermissions []*rmodel.RobotPermission

	// Project permissions
	for projectName, projectPermissions := range b.projects {
		mergedPermissions = append(mergedPermissions, &rmodel.RobotPermission{
			Namespace: projectName,
			Access:    permissionsToAccesses(projectPermissions),
			Kind:      "project",
		})
	}

	// System permissions
	if len(b.system) > 0 {
		mergedPermissions = append(mergedPermissions, &rmodel.RobotPermission{
			Namespace: "/",
			Access:    permissionsToAccesses(b.system),
			Kind:      "system",
		})
	}
	return mergedPermissions
}

// permissionsToAccesses converts model permissions into Harbor API Access entries.
func permissionsToAccesses(perms []models.Permission) []*models.Access {
	var acc []*models.Access
	for _, p := range perms {
		acc = append(acc, &models.Access{
			Resource: p.Resource,
			Action:   p.Action,
		})
	}
	return acc
}
