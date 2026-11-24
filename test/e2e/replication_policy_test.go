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

//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/goharbor/harbor-cli/cmd/harbor/root/replication/policies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type replicationFilter struct {
	Type       string `json:"type"`
	Decoration string `json:"decoration"`
	Value      any    `json:"value"`
}

type replicationPolicy struct {
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Filters     []replicationFilter `json:"filters"`
}

// findIDByName looks up a resource created via POST, since Harbor returns
// only a Location header there.
func findIDByName(t *testing.T, path, name string) int64 {
	t.Helper()

	var items []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	doJSON(t, http.MethodGet, path+"?q="+url.QueryEscape("name="+name), nil, &items)
	for _, it := range items {
		if it.Name == name {
			return it.ID
		}
	}
	t.Fatalf("%s: no item named %q", path, name)
	return 0
}

// createPushPolicy creates a registry endpoint and a push replication policy
// with the given filters, and registers both for removal. Harbor health-checks
// the endpoint on create, so it points at RegistryURL.
func createPushPolicy(t *testing.T, prefix string, filters []replicationFilter) int64 {
	t.Helper()

	suffix := strconv.FormatInt(time.Now().UnixNano(), 10)
	registryName := prefix + "-reg-" + suffix
	policyName := prefix + "-policy-" + suffix

	doJSON(t, http.MethodPost, "/api/v2.0/registries", map[string]any{
		"name":     registryName,
		"type":     "docker-registry",
		"url":      RegistryURL(),
		"insecure": true,
	}, nil)
	registryID := findIDByName(t, "/api/v2.0/registries", registryName)
	t.Cleanup(func() {
		doJSON(t, http.MethodDelete, fmt.Sprintf("/api/v2.0/registries/%d", registryID), nil, nil)
	})

	doJSON(t, http.MethodPost, "/api/v2.0/replication/policies", map[string]any{
		"name":          policyName,
		"description":   "before",
		"enabled":       false,
		"dest_registry": map[string]any{"id": registryID},
		"trigger":       map[string]any{"type": "manual"},
		"filters":       filters,
	}, nil)
	policyID := findIDByName(t, "/api/v2.0/replication/policies", policyName)
	t.Cleanup(func() {
		doJSON(t, http.MethodDelete, fmt.Sprintf("/api/v2.0/replication/policies/%d", policyID), nil, nil)
	})

	return policyID
}

func getPolicy(t *testing.T, id int64) replicationPolicy {
	t.Helper()

	var p replicationPolicy
	doJSON(t, http.MethodGet, fmt.Sprintf("/api/v2.0/replication/policies/%d", id), nil, &p)
	return p
}

func labelFilter(t *testing.T, p replicationPolicy) replicationFilter {
	t.Helper()

	for _, f := range p.Filters {
		if f.Type == "label" {
			return f
		}
	}
	t.Fatalf("policy %d has no label filter: %+v", p.ID, p.Filters)
	return replicationFilter{}
}

// Test_ReplicationPolicyUpdate_DescriptionKeepsLabelFilter covers #1129: a
// description-only update must not split a label name that contains a comma.
func Test_ReplicationPolicyUpdate_DescriptionKeepsLabelFilter(t *testing.T) {
	for _, tc := range []struct {
		name   string
		labels []string
	}{
		{name: "separate labels", labels: []string{"release", "approved"}},
		{name: "label containing comma", labels: []string{"release,approved"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			policyID := createPushPolicy(t, "e2e-1129", []replicationFilter{
				{Type: "label", Decoration: "matches", Value: tc.labels},
			})
			before := labelFilter(t, getPolicy(t, policyID))

			loginAndReload(t)

			update := policies.UpdateCommand()
			update.SetArgs([]string{strconv.FormatInt(policyID, 10), "--description", "after"})
			require.NoError(t, update.Execute())

			after := getPolicy(t, policyID)
			assert.Equal(t, "after", after.Description)
			assert.Equal(t, before, labelFilter(t, after), "description-only update changed the label filter")
		})
	}
}
