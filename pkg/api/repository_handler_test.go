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
package api

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// useTestServer points the shared Harbor client at a local server that
// records the escaped request path and answers with the given body. The
// returned function reads the last recorded path.
func useTestServer(t *testing.T, body string) func() string {
	t.Helper()

	var mu sync.Mutex
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gotPath = r.URL.EscapedPath()
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))

	utils.ClientOnce = sync.Once{}
	utils.ClientOnce.Do(func() {})
	utils.ClientInstance = utils.GetClientByConfig(&harbor.ClientSetConfig{URL: srv.URL})
	utils.ClientErr = nil

	t.Cleanup(func() {
		srv.Close()
		utils.ClientOnce = sync.Once{}
		utils.ClientInstance = nil
		utils.ClientErr = nil
	})

	return func() string {
		mu.Lock()
		defer mu.Unlock()
		return gotPath
	}
}

// Harbor expects a slash in repository_name to be encoded twice
// (a/b -> a%2Fb -> a%252Fb), otherwise the request does not reach the
// repository handler when Harbor runs behind its nginx proxy.
func TestRepositoryNameIsEncodedInPath(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		call     func() error
		wantPath string
	}{
		{
			name: "list artifacts of nested repository",
			body: "[]",
			call: func() error {
				_, err := ListArtifact("library", "team/app")
				return err
			},
			wantPath: "/api/v2.0/projects/library/repositories/team%252Fapp/artifacts",
		},
		{
			name: "view artifact of nested repository",
			body: "{}",
			call: func() error {
				_, err := ViewArtifact("library", "team/app", "latest", false)
				return err
			},
			wantPath: "/api/v2.0/projects/library/repositories/team%252Fapp/artifacts/latest",
		},
		{
			name: "list tags of nested repository",
			body: "[]",
			call: func() error {
				_, err := ListTags("library", "team/app", "latest")
				return err
			},
			wantPath: "/api/v2.0/projects/library/repositories/team%252Fapp/artifacts/latest/tags",
		},
		{
			name: "view nested repository",
			body: "{}",
			call: func() error {
				_, err := RepoView("library", "team/app")
				return err
			},
			wantPath: "/api/v2.0/projects/library/repositories/team%252Fapp",
		},
		{
			name: "list artifacts of single-level repository",
			body: "[]",
			call: func() error {
				_, err := ListArtifact("library", "nginx")
				return err
			},
			wantPath: "/api/v2.0/projects/library/repositories/nginx/artifacts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath := useTestServer(t, tt.body)

			assert.NoError(t, tt.call())
			assert.Equal(t, tt.wantPath, gotPath())
		})
	}
}
