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
package replication

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/goharbor/go-client/pkg/harbor"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

const childEnv = "HARBOR_CLI_REPLICATION_PROMPT_CHILD"

// useReplicationServer fakes Harbor: listing policies fails with a 500,
// while executions and tasks come back as empty lists.
func useReplicationServer(t *testing.T) {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2.0/replication/policies":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"errors":[{"code":"UNKNOWN","message":"boom"}]}`))
		case "/api/v2.0/replication/executions", "/api/v2.0/replication/executions/7/tasks":
			_, _ = w.Write([]byte(`[]`))
		default:
			http.NotFound(w, r)
		}
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
}

// The prompt helpers used to call log.Fatal, which kills the process. Each
// case runs in a child process so a hard exit shows up as a test failure
// instead of taking the whole test binary down.
func TestReplicationPromptsReturnErrors(t *testing.T) {
	tests := []struct {
		name    string
		cmd     func() *cobra.Command
		args    []string
		wantErr string
	}{
		// A nil slice makes cobra fall back to os.Args, so pass an empty one.
		{"start_policy_list_fails", StartCommand, []string{}, "failed to select replication policy: boom"},
		{"stop_policy_list_fails", StopCommand, []string{}, "failed to select replication policy: boom"},
		{"stop_no_executions", StopCommand, []string{"5"}, "failed to select replication execution: no replication executions found"},
		{"log_policy_list_fails", LogsCommand, []string{}, "failed to select replication policy: boom"},
		{"log_no_tasks", LogsCommand, []string{"--execution-id", "7"}, "failed to select replication task: no replication tasks found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv(childEnv) == t.Name() {
				useReplicationServer(t)
				cmd := tt.cmd()
				cmd.SetArgs(tt.args)
				cmd.SilenceUsage = true
				require.EqualError(t, cmd.Execute(), tt.wantErr)
				return
			}

			child := exec.Command(os.Args[0], "-test.run=^"+t.Name()+"$", "-test.v")
			child.Env = append(os.Environ(), childEnv+"="+t.Name())
			out, err := child.CombinedOutput()
			require.NoError(t, err, "command exited instead of returning an error:\n%s", out)
		})
	}
}
