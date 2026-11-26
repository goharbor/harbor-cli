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
	"errors"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/jobservice"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
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

func TestGetWorkerPools_Success(t *testing.T) {
	payload := []*models.WorkerPool{
		{WorkerPoolID: "pool-abc123", Host: "10.0.0.1:8080", Pid: 12345, Concurrency: 10},
	}
	injectJobserviceClient(t, &jobservice.GetWorkerPoolsOK{Payload: payload}, nil)

	response, err := GetWorkerPools()

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, payload, response.Payload)
}

func TestGetWorkerPools_CallError(t *testing.T) {
	injectJobserviceClient(t, nil, errors.New("boom"))

	response, err := GetWorkerPools()

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "boom")
}

func TestGetWorkers_Success(t *testing.T) {
	payload := []*models.Worker{
		{ID: "worker-1", PoolID: "pool-abc123", JobID: "job-9", JobName: "REPLICATE"},
	}
	injectJobserviceClient(t, &jobservice.GetWorkersOK{Payload: payload}, nil)

	response, err := GetWorkers("pool-abc123")

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, payload, response.Payload)
}

func TestGetWorkers_CallError(t *testing.T) {
	injectJobserviceClient(t, nil, &jobservice.GetWorkersUnauthorized{})

	response, err := GetWorkers("pool-abc123")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "Unauthorized")
}

func TestGetWorkerPools_ClientInitError(t *testing.T) {
	origInstance := utils.ClientInstance
	origErr := utils.ClientErr
	t.Cleanup(func() {
		utils.ClientInstance = origInstance
		utils.ClientErr = origErr
	})

	utils.ClientOnce.Do(func() {})
	utils.ClientInstance = nil
	utils.ClientErr = errors.New("no credentials")

	response, err := GetWorkerPools()
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "no credentials")

	utils.ClientErr = errors.New("no credentials")
	workerResponse, workerErr := GetWorkers("pool-abc123")
	assert.Error(t, workerErr)
	assert.Nil(t, workerResponse)
	assert.Contains(t, workerErr.Error(), "no credentials")
}
