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

type ListFlags struct {
	ProjectID int64
	Scope     string
	Name      string
	Page      int64
	PageSize  int64
	Q         string
	Sort      string
	Public    bool
}

// CreateView for Registry
type CreateRegView struct {
	Name        string
	Type        string
	Description string
	URL         string
	Credential  RegistryCredential
	Insecure    bool
}

type ContextListView struct {
	Name     string
	Username string
	Server   string
}

// Credential for Registry
type RegistryCredential struct {
	AccessKey    string `json:"access_key,omitempty"`
	Type         string `json:"type,omitempty"`
	AccessSecret string `json:"access_secret,omitempty"`
}

type ListQuotaFlags struct {
	PageSize    int64
	Page        int64
	Sort        string
	Reference   string
	ReferenceID string
}

// Provides type for List member
type ListMemberOptions struct {
	XIsResourceName bool
	ProjectNameOrID string
	Page            int64
	PageSize        int64
	EntityName      string
	WithDetail      bool
}

// Provides type for Update member
type UpdateMemberOptions struct {
	XIsResourceName bool
	ID              int64
	ProjectNameOrID string
	RoleID          *models.RoleRequest
}

// Provides Params for getting the Member
type GetMemberOptions struct {
	XIsResourceName bool
	ID              int64
	ProjectNameOrID string
}
