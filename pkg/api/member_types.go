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
package api

import "github.com/goharbor/go-client/pkg/sdk/v2.0/models"

// ListMemberOptions provides options for listing project members.
type ListMemberOptions struct {
	XIsResourceName bool
	ProjectNameOrID string
	Page            int64
	PageSize        int64
	EntityName      string
	WithDetail      bool
}

// UpdateMemberOptions provides options for updating a project member.
type UpdateMemberOptions struct {
	XIsResourceName bool
	ID              int64
	ProjectNameOrID string
	RoleID          *models.RoleRequest
}

// GetMemberOptions provides parameters for getting a specific project member.
type GetMemberOptions struct {
	XIsResourceName bool
	ID              int64
	ProjectNameOrID string
}
