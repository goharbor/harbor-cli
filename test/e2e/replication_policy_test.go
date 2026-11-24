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

// labelFilters returns every label filter of p, in order.
func labelFilters(t *testing.T, p replicationPolicy) []replicationFilter {
	t.Helper()

	var out []replicationFilter
	for _, f := range p.Filters {
		if f.Type == "label" {
			out = append(out, f)
		}
	}
	require.NotEmpty(t, out, "policy %d has no label filter: %+v", p.ID, p.Filters)
	return out
}

// Test_ReplicationPolicyUpdate_KeepsLabelFilters covers #1129: updates that
// do not supply a new label pattern must keep every existing label filter and
// must not split a label name that contains a comma.
func Test_ReplicationPolicyUpdate_KeepsLabelFilters(t *testing.T) {
	matches := func(labels ...string) replicationFilter {
		return replicationFilter{Type: "label", Decoration: "matches", Value: labels}
	}
	for _, tc := range []struct {
		name    string
		filters []replicationFilter
		args    []string
		// decoration, when set, is the expected decoration of every label filter.
		decoration string
	}{
		{name: "separate labels", filters: []replicationFilter{matches("release", "approved")}, args: []string{"--description", "after"}},
		{name: "label containing comma", filters: []replicationFilter{matches("release,approved")}, args: []string{"--description", "after"}},
		{name: "multiple filters", filters: []replicationFilter{matches("release,approved"), matches("beta")}, args: []string{"--description", "after"}},
		{name: "filter type only", filters: []replicationFilter{matches("release,approved"), matches("beta")}, args: []string{"--label-filter", "excludes"}, decoration: "excludes"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			policyID := createPushPolicy(t, "e2e-1129", tc.filters)
			want := labelFilters(t, getPolicy(t, policyID))
			require.Len(t, want, len(tc.filters))
			if tc.decoration != "" {
				for i := range want {
					want[i].Decoration = tc.decoration
				}
			}

			loginAndReload(t)

			update := policies.UpdateCommand()
			update.SetArgs(append([]string{strconv.FormatInt(policyID, 10)}, tc.args...))
			require.NoError(t, update.Execute())

			after := getPolicy(t, policyID)
			if tc.args[0] == "--description" {
				assert.Equal(t, tc.args[1], after.Description)
			}
			assert.Equal(t, want, labelFilters(t, after), "update changed the label filters")
		})
	}
}
