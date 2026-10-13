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
	"os"
	"path/filepath"
	"testing"

	"github.com/goharbor/harbor-cli/pkg/declarative"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeOverlaysNamedResourcesAndSpecifiedFields(t *testing.T) {
	autoScan := true
	enabled := true
	base := declarative.NewConfiguration()
	base.Spec.System = map[string]any{"read_only": false, "self_registration": true}
	base.Spec.Projects = []declarative.Project{{
		Name:     "application",
		Metadata: &declarative.ProjectMetadata{AutoScan: &autoScan},
		Quota:    map[string]int64{"storage": 100, "count": -1},
	}}
	base.Spec.ReplicationPolicies = []declarative.ReplicationPolicy{{
		Name: "mirror", Mode: "push", Registry: "hub",
		Trigger: &declarative.ReplicationTrigger{Type: "scheduled", Cron: "0 0 * * * *"},
	}}

	overlay := declarative.NewConfiguration()
	overlay.Spec.System = map[string]any{"read_only": true}
	overlay.Spec.Projects = []declarative.Project{{
		Name:     "application",
		Metadata: &declarative.ProjectMetadata{Severity: ptr("critical")},
		Quota:    map[string]int64{"storage": 500},
	}}
	overlay.Spec.ReplicationPolicies = []declarative.ReplicationPolicy{{
		Name: "mirror", Enabled: &enabled,
		Trigger: &declarative.ReplicationTrigger{Type: "manual"},
	}}

	merged, err := declarative.Merge(base, overlay)
	require.NoError(t, err)

	assert.Equal(t, true, merged.Spec.System["read_only"])
	assert.Equal(t, true, merged.Spec.System["self_registration"])
	require.Len(t, merged.Spec.Projects, 1)
	assert.Equal(t, int64(500), merged.Spec.Projects[0].Quota["storage"])
	assert.Equal(t, int64(-1), merged.Spec.Projects[0].Quota["count"])
	assert.Equal(t, &autoScan, merged.Spec.Projects[0].Metadata.AutoScan)
	assert.Equal(t, "critical", *merged.Spec.Projects[0].Metadata.Severity)
	require.Len(t, merged.Spec.ReplicationPolicies, 1)
	assert.Equal(t, "push", merged.Spec.ReplicationPolicies[0].Mode)
	assert.Equal(t, "hub", merged.Spec.ReplicationPolicies[0].Registry)
	assert.Equal(t, &enabled, merged.Spec.ReplicationPolicies[0].Enabled)
	assert.Equal(t, "manual", merged.Spec.ReplicationPolicies[0].Trigger.Type)
	assert.Empty(t, merged.Spec.ReplicationPolicies[0].Trigger.Cron)
	assert.Equal(t, int64(100), base.Spec.Projects[0].Quota["storage"], "merge must not mutate its inputs")
}

func TestReadPathRecursivelyAppliesFilesInLexicalOrder(t *testing.T) {
	directory := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(directory, "projects"), 0o755))
	require.NoError(t, os.Mkdir(filepath.Join(directory, ".ignored"), 0o755))
	writeTestFile(t, filepath.Join(directory, "00-base.yaml"), `apiVersion: goharbor.io/v1alpha1
kind: HarborConfiguration
spec:
  projects:
    - name: application
      public: false
      quota:
        storage: 100
  replicationPolicies:
    - name: mirror
      mode: push
      registry: hub
`)
	writeTestFile(t, filepath.Join(directory, "10-production.yaml"), `apiVersion: goharbor.io/v1alpha1
kind: HarborConfiguration
spec:
  projects:
    - name: application
      quota:
        storage: 500
  replicationPolicies:
    - name: mirror
      enabled: true
`)
	writeTestFile(t, filepath.Join(directory, "projects", "20-metadata.json"), `{
  "apiVersion": "goharbor.io/v1alpha1",
  "kind": "HarborConfiguration",
  "spec": {"projects": [{"name": "application", "metadata": {"severity": "critical"}}]}
}`)
	writeTestFile(t, filepath.Join(directory, "README.md"), "not configuration")
	writeTestFile(t, filepath.Join(directory, ".ignored", "invalid.yaml"), "not: a configuration")

	configuration, err := declarative.ReadPath(directory)
	require.NoError(t, err)

	require.Len(t, configuration.Spec.Projects, 1)
	project := configuration.Spec.Projects[0]
	assert.Equal(t, int64(500), project.Quota["storage"])
	assert.Equal(t, "critical", *project.Metadata.Severity)
	require.Len(t, configuration.Spec.ReplicationPolicies, 1)
	assert.Equal(t, "push", configuration.Spec.ReplicationPolicies[0].Mode)
	assert.True(t, *configuration.Spec.ReplicationPolicies[0].Enabled)
}

func TestReadPathRejectsEmptyDirectory(t *testing.T) {
	_, err := declarative.ReadPath(t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "contains no YAML or JSON files")
}

func writeTestFile(t *testing.T, filename, contents string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filename, []byte(contents), 0o600))
}
