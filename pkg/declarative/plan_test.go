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

package declarative_test

import (
	"testing"

	"github.com/goharbor/harbor-cli/pkg/declarative"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildPlanOrdersDependenciesAndIgnoresUnmanagedFields(t *testing.T) {
	public := false
	registryType := "docker-hub"
	desired := declarative.NewConfiguration()
	desired.Spec.Registries = []declarative.Registry{{Name: "hub", Type: &registryType}}
	desired.Spec.Projects = []declarative.Project{
		{
			Name:   "demo",
			Public: &public,
			Quota:  map[string]int64{"storage": 100},
			Webhooks: []declarative.Webhook{
				{Name: "notify"},
			},
		},
	}
	desired.Spec.ReplicationPolicies = []declarative.ReplicationPolicy{{Name: "mirror", Mode: "push", Registry: "hub"}}
	desired.Spec.System = map[string]any{"read_only": false}

	current := declarative.NewConfiguration()
	current.Spec.Registries = []declarative.Registry{{Name: "hub", Type: &registryType, URL: ptr("https://hub.docker.com")}}
	current.Spec.Projects = []declarative.Project{{Name: "demo", Public: &public, Quota: map[string]int64{"storage": 50}}}

	plan, err := declarative.BuildPlan(desired, current)
	require.NoError(t, err)

	operations := make([]string, 0, len(plan.Actions))
	for _, action := range plan.Actions {
		operations = append(operations, action.String())
	}
	assert.Equal(t, []string{
		"no-op registry hub",
		"no-op project demo",
		"update quota demo",
		"create webhook demo/notify",
		"create replication policy mirror",
		"update system configuration",
	}, operations)
	assert.Equal(t, 4, plan.ChangeCount())
}

func TestBuildPlanIsNoopAfterConvergence(t *testing.T) {
	enabled := true
	document := declarative.NewConfiguration()
	document.Spec.ReplicationPolicies = []declarative.ReplicationPolicy{
		{Name: "mirror", Mode: "pull", Registry: "source", Enabled: &enabled},
	}

	plan, err := declarative.BuildPlan(document, document)
	require.NoError(t, err)
	require.Len(t, plan.Actions, 1)
	assert.Equal(t, declarative.OperationNoop, plan.Actions[0].Operation)
	assert.False(t, plan.HasChanges())
}

func ptr[T any](value T) *T {
	return &value
}
