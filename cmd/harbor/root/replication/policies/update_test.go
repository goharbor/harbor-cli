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
package policies

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// labelFilter returns the JSON of a matches label filter with the given value.
func labelFilter(value string) string {
	return `{"type":"label","decoration":"matches","value":` + value + `}`
}

// runUpdate executes the update command against an in-memory Harbor that
// serves a policy with the given label value and returns the filters sent in the PUT.
func runUpdate(t *testing.T, labels string, args ...string) []map[string]any {
	t.Helper()
	return runUpdateWithFilters(t, "["+labelFilter(labels)+"]", args...)
}

// runUpdateWithFilters is runUpdate for a policy with the given JSON filters array.
func runUpdateWithFilters(t *testing.T, filters string, args ...string) []map[string]any {
	t.Helper()

	policy := `{"id":42,"name":"p","description":"before","dest_registry":{"id":9},"trigger":{"type":"manual"},` +
		`"filters":` + filters + `}`

	var sent struct {
		Filters []map[string]any `json:"filters"`
	}
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		body := policy
		if r.Method == http.MethodPut {
			if err := json.NewDecoder(r.Body).Decode(&sent); err != nil {
				return nil, err
			}
			body = ""
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    r,
		}, nil
	})

	utils.ClientOnce = sync.Once{}
	utils.ClientOnce.Do(func() {})
	utils.ClientInstance = v2client.New(v2client.Config{Transport: transport})
	utils.ClientErr = nil
	t.Cleanup(func() {
		utils.ClientOnce = sync.Once{}
		utils.ClientInstance = nil
		utils.ClientErr = nil
	})

	cmd := UpdateCommand()
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs(append([]string{"42"}, args...))
	require.NoError(t, cmd.Execute())
	return sent.Filters
}

func TestUpdateKeepsLabelFilterWhenNotEdited(t *testing.T) {
	for _, labels := range []string{`["release","approved"]`, `["release,approved"]`} {
		t.Run(labels, func(t *testing.T) {
			filters := runUpdate(t, labels, "--description", "after")
			require.Len(t, filters, 1)

			var want []any
			require.NoError(t, json.Unmarshal([]byte(labels), &want))
			assert.Equal(t, want, filters[0]["value"])
		})
	}
}

func TestUpdateRebuildsLabelFilterWhenEdited(t *testing.T) {
	filters := runUpdate(t, `["old"]`, "--label-pattern", "a, b")
	require.Len(t, filters, 1)
	assert.Equal(t, []any{"a", "b"}, filters[0]["value"])
}

func TestUpdateKeepsAllLabelFiltersWhenNotEdited(t *testing.T) {
	filters := runUpdateWithFilters(t,
		"["+labelFilter(`["release"]`)+","+labelFilter(`["approved"]`)+"]",
		"--description", "after")
	require.Len(t, filters, 2)
	assert.Equal(t, []any{"release"}, filters[0]["value"])
	assert.Equal(t, []any{"approved"}, filters[1]["value"])
}

func TestUpdateRebuildsLabelFilterWhenPatternMatchesBaseline(t *testing.T) {
	// --label-pattern uses comma-separated syntax, so "release,approved" asks
	// for two labels even though it equals the joined existing value.
	filters := runUpdate(t, `["release,approved"]`, "--label-pattern", "release,approved")
	require.Len(t, filters, 1)
	assert.Equal(t, []any{"release", "approved"}, filters[0]["value"])
}

func TestUpdateLabelFilterOnlyKeepsLabelValues(t *testing.T) {
	filters := runUpdateWithFilters(t,
		"["+labelFilter(`["release,approved"]`)+","+labelFilter(`["beta"]`)+"]",
		"--label-filter", "excludes")
	require.Len(t, filters, 2)
	assert.Equal(t, []any{"release,approved"}, filters[0]["value"])
	assert.Equal(t, []any{"beta"}, filters[1]["value"])
	for _, f := range filters {
		assert.Equal(t, "excludes", f["decoration"])
	}
}

func TestUpdateLabelFilterAppliesToEveryFilter(t *testing.T) {
	excludes := `{"type":"label","decoration":"excludes","value":["beta"]}`
	filters := runUpdateWithFilters(t,
		"["+labelFilter(`["release"]`)+","+excludes+"]",
		"--label-filter", "excludes")
	require.Len(t, filters, 2)
	for _, f := range filters {
		assert.Equal(t, "excludes", f["decoration"])
	}
}
