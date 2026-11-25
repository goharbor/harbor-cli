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
	"unicode"

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

// RobotPermission is the SDK permission scope; kept as an alias so view
// consumers do not need a conversion layer.
type RobotPermission = models.RobotPermission

// CreateView carries the user's input for a new robot account.
type CreateView struct {
	Description string             `json:"description,omitempty"`
	Disable     bool               `json:"disable,omitempty"`
	Duration    int64              `json:"duration,omitempty"`
	Level       string             `json:"level,omitempty"`
	Name        string             `json:"name,omitempty"`
	Permissions []*RobotPermission `json:"permissions"`
	Secret      string             `json:"secret,omitempty"`
	ProjectName string
}

// UpdateView carries the full state of an existing robot account while it
// is being edited.
type UpdateView struct {
	CreationTime strfmt.DateTime    `json:"creation_time,omitempty"`
	Description  string             `json:"description,omitempty"`
	Disable      bool               `json:"disable,omitempty"`
	Duration     int64              `json:"duration,omitempty"`
	Editable     bool               `json:"editable"`
	ExpiresAt    int64              `json:"expires_at,omitempty"`
	ID           int64              `json:"id,omitempty"`
	Level        string             `json:"level,omitempty"`
	Name         string             `json:"name,omitempty"`
	Permissions  []*RobotPermission `json:"permissions"`
	Secret       string             `json:"secret,omitempty"`
	UpdateTime   strfmt.DateTime    `json:"update_time,omitempty"`
}

func isValidName(s string) bool {
	if s == "" {
		return false
	}
	if !unicode.IsLetter(rune(s[0])) {
		return false
	}
	for _, r := range s {
		if unicode.IsUpper(r) {
			return false
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false
		}
	}
	return true
}
