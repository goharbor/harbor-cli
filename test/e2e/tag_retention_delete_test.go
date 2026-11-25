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
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// createRetentionPolicy gives project a tag retention policy with one rule.
func createRetentionPolicy(t *testing.T, project string) {
	t.Helper()

	var p struct {
		ProjectID int64 `json:"project_id"`
	}
	doJSON(t, http.MethodGet, "/api/v2.0/projects/"+url.PathEscape(project), nil, &p)

	doJSON(t, http.MethodPost, "/api/v2.0/retentions", map[string]any{
		"algorithm": "or",
		"rules": []map[string]any{{
			"action":   "retain",
			"template": "always",
			"params":   map[string]any{},
			"tag_selectors": []map[string]string{
				{"kind": "doublestar", "decoration": "matches", "pattern": "**"},
			},
			"scope_selectors": map[string]any{
				"repository": []map[string]string{
					{"kind": "doublestar", "decoration": "repoMatches", "pattern": "**"},
				},
			},
		}},
		"trigger": map[string]any{"kind": "Schedule", "settings": map[string]string{"cron": ""}},
		"scope":   map[string]any{"level": "project", "ref": p.ProjectID},
	}, nil)
}

// createGuest creates a user with guest access to project and registers its
// removal with t.Cleanup. Guests can read the project but not its retention
// policy.
func createGuest(t *testing.T, project, username, password string) {
	t.Helper()

	doJSON(t, http.MethodPost, "/api/v2.0/users", map[string]any{
		"username": username,
		"email":    username + "@example.com",
		"realname": username,
		"password": password,
	}, nil)

	t.Cleanup(func() {
		var users []struct {
			UserID int64 `json:"user_id"`
		}
		doJSON(t, http.MethodGet, "/api/v2.0/users/search?username="+url.QueryEscape(username), nil, &users)
		for _, u := range users {
			doJSON(t, http.MethodDelete, fmt.Sprintf("/api/v2.0/users/%d", u.UserID), nil, nil)
		}
	})

	doJSON(t, http.MethodPost, "/api/v2.0/projects/"+url.PathEscape(project)+"/members", map[string]any{
		"role_id":     3, // guest
		"member_user": map[string]string{"username": username},
	}, nil)
}

// A guest can resolve the project's retention ID but is not allowed to read
// the policy. The failed rule lookup must come back as an error instead of
// crashing the CLI.
func Test_TagRetentionDelete_ReturnsRuleLookupError(t *testing.T) {
	suffix := time.Now().UnixNano()
	project := fmt.Sprintf("e2e-995-%d", suffix)
	guest := fmt.Sprintf("e2e-995-guest-%d", suffix)
	const guestPassword = "Harbor12345x"

	createProject(t, project)
	createRetentionPolicy(t, project)
	createGuest(t, project, guest, guestPassword)

	require.NoError(t, login(t, Server(), guest, guestPassword))

	err := harbor(t, "tag", "retention", "delete", "--project-name", project)
	require.Error(t, err)
	require.ErrorContains(t, err, "failed to get retention rule")
	require.ErrorContains(t, err, "403")
}
