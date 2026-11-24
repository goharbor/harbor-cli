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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"

	"github.com/goharbor/harbor-cli/cmd/harbor/root"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/require"
)

// loginAndReload runs `harbor login` as Username, then reloads the saved
// config the way a new CLI process would, so later commands and direct API
// helper calls pick up the stored credential.
func loginAndReload(t *testing.T) {
	t.Helper()

	require.NoError(t, login(t, Server(), Username(), Password()))
	reloadConfig(t)
}

// reloadConfig drops the cached config and API client and loads the config
// file that `harbor login` wrote. It returns that file's path.
func reloadConfig(t *testing.T) string {
	t.Helper()

	data, err := utils.GetCurrentHarborData()
	require.NoError(t, err)

	utils.ConfigInitialization = &utils.Once{}
	utils.CurrentHarborConfig = nil
	utils.ClientOnce = sync.Once{}
	utils.ClientInstance = nil
	utils.ClientErr = nil
	utils.InitConfig(data.ConfigPath, true)

	return data.ConfigPath
}

// harbor runs the CLI with args as a fresh invocation would, against the
// config that `harbor login` wrote earlier in the test.
func harbor(t *testing.T, args ...string) error {
	t.Helper()

	configPath := reloadConfig(t)

	cmd := root.RootCmd()
	cmd.SetArgs(append([]string{"--config", configPath}, args...))
	return cmd.Execute()
}

// createProject creates a private project and registers its removal with
// t.Cleanup.
func createProject(t *testing.T, name string) {
	t.Helper()

	doJSON(t, http.MethodPost, "/api/v2.0/projects", map[string]any{
		"project_name": name,
		"metadata":     map[string]string{"public": "false"},
	}, nil)

	deleteProjectOnCleanup(t, name)
}

// deleteProjectOnCleanup removes project when the test ends. A missing
// project is fine, since some tests expect the CLI never to create it.
func deleteProjectOnCleanup(t *testing.T, name string) {
	t.Helper()

	t.Cleanup(func() {
		status := projectStatus(t, http.MethodDelete, name)
		if status != http.StatusOK && status != http.StatusNotFound {
			t.Errorf("delete project %s: unexpected status %d", name, status)
		}
	})
}

// projectStatus sends method to the project's API endpoint and returns the
// response status code.
func projectStatus(t *testing.T, method, name string) int {
	t.Helper()

	req, err := http.NewRequest(method, Server()+"/api/v2.0/projects/"+url.PathEscape(name), nil)
	require.NoError(t, err)
	req.SetBasicAuth(Username(), Password())

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	require.NoError(t, err)
	_ = resp.Body.Close()

	return resp.StatusCode
}

// captureStdout runs fn and returns what it wrote to os.Stdout.
func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = w
	// Buffered so the reader never blocks if we return early.
	done := make(chan []byte, 1)
	go func() {
		out, _ := io.ReadAll(r)
		done <- out
	}()

	// Restore stdout and close the pipe even if fn calls t.FailNow.
	var restoreOnce sync.Once
	restore := func() {
		restoreOnce.Do(func() {
			os.Stdout = orig
			_ = w.Close()
		})
	}
	defer func() {
		restore()
		_ = r.Close()
	}()

	runErr := fn()
	restore()
	out := <-done

	return string(bytes.TrimSpace(out)), runErr
}

// createUser creates a local Harbor user and returns its ID. The user is
// removed with t.Cleanup unless the test already deleted it.
func createUser(t *testing.T, name string) int64 {
	t.Helper()

	doJSON(t, http.MethodPost, "/api/v2.0/users", map[string]string{
		"username": name,
		"email":    emailLocalPart(name) + "@example.com",
		"realname": name,
		"password": "Harbor12345",
	}, nil)
	// Registered before the lookup so a failed lookup does not leak the user.
	t.Cleanup(func() {
		for _, id := range findUserIDs(t, name) {
			doJSON(t, http.MethodDelete, fmt.Sprintf("/api/v2.0/users/%d", id), nil, nil)
		}
	})

	ids := findUserIDs(t, name)
	require.Len(t, ids, 1, "user %s not found after create", name)
	return ids[0]
}

// emailLocalPart drops characters from name that are not valid in an email
// address.
func emailLocalPart(name string) string {
	return strings.Map(func(r rune) rune {
		if r == '-' || r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, name)
}

// findUserIDs returns the IDs of users named exactly name. It uses a
// substring query because Harbor rejects exact values that look like ranges,
// and escapes name twice because Harbor URL-decodes q once more.
func findUserIDs(t *testing.T, name string) []int64 {
	t.Helper()

	var found []struct {
		UserID   int64  `json:"user_id"`
		Username string `json:"username"`
	}
	doJSON(t, http.MethodGet, "/api/v2.0/users?page_size=100&q="+url.QueryEscape("username=~"+url.QueryEscape(name)), nil, &found)
	var ids []int64
	for _, u := range found {
		if u.Username == name {
			ids = append(ids, u.UserID)
		}
	}
	return ids
}

// userExists reports whether Harbor still has a user with id.
func userExists(t *testing.T, id int64) bool {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2.0/users/%d", Server(), id), nil)
	require.NoError(t, err)
	req.SetBasicAuth(Username(), Password())

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true
	case http.StatusNotFound:
		return false
	default:
		t.Fatalf("GET user %d: unexpected status %s", id, resp.Status)
		return false
	}
}
