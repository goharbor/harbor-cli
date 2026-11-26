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
package workers

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/jobservice"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type stubJobserviceTransport struct {
	result interface{}
	err    error
}

func (s *stubJobserviceTransport) Submit(op *runtime.ClientOperation) (interface{}, error) {
	return s.result, s.err
}

func injectJobserviceClient(t *testing.T, result interface{}, callErr error) {
	t.Helper()

	origInstance := utils.ClientInstance
	origErr := utils.ClientErr
	t.Cleanup(func() {
		utils.ClientInstance = origInstance
		utils.ClientErr = origErr
	})

	utils.ClientInstance = &v2client.HarborAPI{
		Jobservice: jobservice.New(&stubJobserviceTransport{result: result, err: callErr}, strfmt.Default, nil),
	}
	utils.ClientErr = nil
	utils.ClientOnce.Do(func() {})
}

func executeListCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()

	var outBuf, cmdBuf bytes.Buffer
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stdout = w

	cmd := ListCommand()
	cmd.SetOut(&cmdBuf)
	cmd.SetErr(&cmdBuf)
	cmd.SetArgs(args)
	cmdErr := cmd.Execute()

	w.Close()
	os.Stdout = oldStdout
	_, _ = io.Copy(&outBuf, r)

	return outBuf.String(), cmdErr
}

func TestWorkersCommand_WiresListSubcommand(t *testing.T) {
	cmd := WorkersCommand()

	assert.Equal(t, "workers", cmd.Use)
	found := false
	for _, sub := range cmd.Commands() {
		if sub.Name() == "list" {
			found = true
			break
		}
	}
	assert.True(t, found, "workers command should wire the list subcommand")
}

func TestListCommand_RequiresOneArg(t *testing.T) {
	_, err := executeListCommand(t)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accepts 1 arg(s), received 0")
}

func TestListCommand_EmptyPoolID(t *testing.T) {
	_, err := executeListCommand(t, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pool-id is required")
}

func TestListCommand_APIError(t *testing.T) {
	injectJobserviceClient(t, nil, errors.New("boom"))

	out, err := executeListCommand(t, "pool-abc123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to retrieve workers")
	assert.Contains(t, err.Error(), "boom")
	assert.NotContains(t, out, "pool-")
}

func TestListCommand_EmptyPayload(t *testing.T) {
	injectJobserviceClient(t, &jobservice.GetWorkersOK{Payload: []*models.Worker{}}, nil)

	out, err := executeListCommand(t, "pool-abc123")

	assert.NoError(t, err)
	assert.Contains(t, out, "No workers found.")
}

func TestListCommand_TableOutput(t *testing.T) {
	viper.Set("output-format", "")
	t.Cleanup(func() { viper.Set("output-format", "") })

	payload := []*models.Worker{
		{ID: "worker-1", PoolID: "pool-abc123", JobID: "job-9", JobName: "REPLICATE"},
	}
	injectJobserviceClient(t, &jobservice.GetWorkersOK{Payload: payload}, nil)

	out, err := executeListCommand(t, "pool-abc123")

	assert.NoError(t, err)
	assert.Contains(t, out, "WORKER_ID")
	assert.Contains(t, out, "worker-1")
	assert.Contains(t, out, "Total: 1 worker(s)")
}

func TestListCommand_JSONFormat(t *testing.T) {
	viper.Set("output-format", "json")
	t.Cleanup(func() { viper.Set("output-format", "") })

	payload := []*models.Worker{
		{ID: "worker-json9", PoolID: "pool-abc123", JobID: "job-9", JobName: "REPLICATE"},
	}
	injectJobserviceClient(t, &jobservice.GetWorkersOK{Payload: payload}, nil)

	out, err := executeListCommand(t, "pool-abc123")

	assert.NoError(t, err)
	assert.Contains(t, out, "worker-json9")
	assert.Contains(t, out, "{")
}
