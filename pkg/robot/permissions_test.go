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
	"errors"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	robotview "github.com/goharbor/harbor-cli/pkg/views/robot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func scopes() []*robotview.RobotPermission {
	return []*robotview.RobotPermission{
		{
			Kind:      "system",
			Namespace: "/",
			Access:    []*models.Access{{Resource: "registry", Action: "list"}},
		},
		{
			Kind:      "project",
			Namespace: "library",
			Access:    []*models.Access{{Resource: "repository", Action: "pull"}},
		},
	}
}

func TestPermissionSetSplitAndMergeRoundTrip(t *testing.T) {
	permissions := NewPermissionSet()
	permissions.Split(scopes())

	require.Len(t, permissions.System, 1)
	require.Len(t, permissions.Project, 1)
	assert.Equal(t, "pull", permissions.Project["library"][0].Action)

	merged := permissions.Merge()
	require.Len(t, merged, 2)
	byKind := map[string]*robotview.RobotPermission{}
	for _, scope := range merged {
		byKind[scope.Kind] = scope
	}
	assert.Equal(t, "/", byKind["system"].Namespace)
	assert.Equal(t, "library", byKind["project"].Namespace)
}

func TestPermissionSetSplitAppendsDuplicateProjectScopes(t *testing.T) {
	permissions := NewPermissionSet()
	permissions.Split([]*robotview.RobotPermission{
		{Kind: "project", Namespace: "library", Access: []*models.Access{{Resource: "repository", Action: "pull"}}},
		{Kind: "project", Namespace: "library", Access: []*models.Access{{Resource: "repository", Action: "push"}}},
	})
	assert.Len(t, permissions.Project["library"], 2)
}

func TestPermissionSetMergeDropsEmptyScopes(t *testing.T) {
	permissions := NewPermissionSet()
	permissions.Project["empty"] = nil

	assert.Empty(t, permissions.Merge())
	assert.True(t, permissions.IsEmpty())
}

func TestViewFromRobotGuardsOptionalFields(t *testing.T) {
	// Regression: a robot without a duration or permissions used to panic.
	view := viewFromRobot(&models.Robot{ID: 7, Name: "robot$ci"})
	assert.Equal(t, int64(0), view.Duration)

	duration := int64(30)
	view = viewFromRobot(&models.Robot{ID: 7, Duration: &duration})
	assert.Equal(t, int64(30), view.Duration)
}

func TestWrapError(t *testing.T) {
	plain := WrapError("list robots", errors.New("boom"))
	assert.ErrorContains(t, plain, "failed to list robots")

	forbidden := WrapError("list robots", errors.New("[GET /robots][403] listRobotForbidden"))
	assert.ErrorContains(t, forbidden, "permission denied: admin privileges are required to list robots")
}

func TestExportSecretToFile(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	require.NoError(t, ExportSecretToFile("ci", "s3cret", "2026-01-01", 123))
	assert.FileExists(t, dir+"/ci-secret.json")
}
