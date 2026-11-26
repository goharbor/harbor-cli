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
package labels

import (
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateLabelCommand_Flags(t *testing.T) {
	cmd := UpdateLableCommand()
	flags := cmd.Flags()

	for _, name := range []string{"name", "color", "description", "project", "project-id", "global"} {
		assert.NotNil(t, flags.Lookup(name), "flag %s should exist", name)
	}
	assert.Nil(t, flags.Lookup("scope"), "update should not expose a scope flag")
}

func TestApplyLabelUpdateFlags(t *testing.T) {
	existing := models.Label{
		Name:        "old-name",
		Color:       "#FFFFFF",
		Description: "old description",
		Scope:       "g",
	}

	tests := []struct {
		name        string
		args        []string
		wantChanged bool
		want        models.Label
	}{
		{
			name:        "no flags keeps existing values",
			args:        []string{},
			wantChanged: false,
			want:        existing,
		},
		{
			name:        "name only",
			args:        []string{"--name", "new-name"},
			wantChanged: true,
			want:        models.Label{Name: "new-name", Color: "#FFFFFF", Description: "old description", Scope: "g"},
		},
		{
			name:        "color only",
			args:        []string{"--color", "#C92100"},
			wantChanged: true,
			want:        models.Label{Name: "old-name", Color: "#C92100", Description: "old description", Scope: "g"},
		},
		{
			name:        "clear description",
			args:        []string{"--description", ""},
			wantChanged: true,
			want:        models.Label{Name: "old-name", Color: "#FFFFFF", Description: "", Scope: "g"},
		},
		{
			name:        "all fields",
			args:        []string{"-n", "new-name", "--color", "#0095D3", "-d", "new description"},
			wantChanged: true,
			want:        models.Label{Name: "new-name", Color: "#0095D3", Description: "new description", Scope: "g"},
		},
		{
			name:        "lookup flags only do not count as updates",
			args:        []string{"--global", "--project-id", "3"},
			wantChanged: false,
			want:        existing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := UpdateLableCommand()
			require.NoError(t, cmd.Flags().Parse(tt.args))

			opts := &models.Label{}
			opts.Name, _ = cmd.Flags().GetString("name")
			opts.Color, _ = cmd.Flags().GetString("color")
			opts.Description, _ = cmd.Flags().GetString("description")

			updateView := existing
			changed := applyLabelUpdateFlags(cmd.Flags(), opts, &updateView)

			assert.Equal(t, tt.wantChanged, changed)
			assert.Equal(t, tt.want, updateView)
		})
	}
}

func TestValidateLabelUpdate(t *testing.T) {
	tests := []struct {
		name    string
		label   models.Label
		wantErr string
	}{
		{"valid", models.Label{Name: "demo", Color: "#C92100"}, ""},
		{"empty description is allowed", models.Label{Name: "demo", Color: "#C92100", Description: ""}, ""},
		{"empty name", models.Label{Name: "", Color: "#C92100"}, "label name cannot be empty"},
		{"empty color", models.Label{Name: "demo", Color: ""}, "label color cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLabelUpdate(&tt.label)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}
