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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createImmutableRule adds a rule to project that protects tags matching
// tagPattern in every repository.
func createImmutableRule(t *testing.T, project, tagPattern string) {
	t.Helper()

	body := map[string]any{
		"tag_selectors": []map[string]string{
			{"kind": "doublestar", "decoration": "matches", "pattern": tagPattern},
		},
		"scope_selectors": map[string]any{
			"repository": []map[string]string{
				{"kind": "doublestar", "decoration": "repoMatches", "pattern": "**"},
			},
		},
	}
	doJSON(t, http.MethodPost, "/api/v2.0/projects/"+url.PathEscape(project)+"/immutabletagrules", body, nil)
}

// Harbor returns 10 immutable rules per page unless asked for more, so a
// project with more than 10 rules checks that the CLI reads every page.
func Test_TagImmutableList_ShowsAllRules(t *testing.T) {
	const ruleCount = 12

	loginAndReload(t)

	project := fmt.Sprintf("e2e-1158-%d", time.Now().UnixNano())
	createProject(t, project)
	for i := range ruleCount {
		createImmutableRule(t, project, fmt.Sprintf("v%d-*", i))
	}

	out, err := captureStdout(t, func() error {
		return harbor(t, "tag", "immutable", "list", project, "--output-format", "json")
	})
	require.NoError(t, err)

	var rules []map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &rules), "output: %s", out)
	require.Len(t, rules, ruleCount)
}
