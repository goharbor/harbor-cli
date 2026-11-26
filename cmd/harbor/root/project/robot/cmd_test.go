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
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommandSurface locks the CLI surface of the project robot commands:
// usage strings and flags must not change unintentionally.
func TestCommandSurface(t *testing.T) {
	tests := []struct {
		use   string
		cmd   *cobra.Command
		flags []string
	}{
		{"create", CreateRobotCommand(), []string{"all-permission", "export-to-file", "project", "name", "description", "duration", "robot-config-file"}},
		{"update [robotID]", UpdateRobotCommand(), []string{"all-permission", "name", "description", "duration", "project", "robot-config-file"}},
		{"delete [robotName]", DeleteRobotCommand(), []string{"project"}},
		{"list [projectName]", ListRobotCommand(), []string{"page", "page-size", "project-id", "query", "sort"}},
		{"view [robotID]", ViewRobotCommand(), []string{"project"}},
		{"refresh [robotID]", RefreshSecretCommand(), []string{"secret", "secret-stdin"}},
	}

	for _, tt := range tests {
		t.Run(tt.use, func(t *testing.T) {
			require.Equal(t, tt.use, tt.cmd.Use)
			for _, flag := range tt.flags {
				assert.NotNil(t, tt.cmd.Flags().Lookup(flag), "missing flag --%s", flag)
			}
		})
	}
}
