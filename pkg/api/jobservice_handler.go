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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/jobservice"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

type JobAction string

const (
	JobActionStop   JobAction = "stop"
	JobActionPause  JobAction = "pause"
	JobActionResume JobAction = "resume"
)

func ListJobQueues() ([]*models.JobQueue, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Jobservice.ListJobQueues(ctx, &jobservice.ListJobQueuesParams{})
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func ActionPendingJobs(jobType string, action JobAction) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Jobservice.ActionPendingJobs(ctx, &jobservice.ActionPendingJobsParams{
		JobType:       jobType,
		ActionRequest: &models.ActionRequest{Action: string(action)},
	})
	if err != nil {
		return err
	}

	return nil
}

func StopJob(jobID string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Jobservice.StopRunningJob(ctx, &jobservice.StopRunningJobParams{
		JobID: jobID,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetJobLog(jobID string) (string, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return "", err
	}

	response, err := client.Jobservice.ActionGetJobLog(ctx, &jobservice.ActionGetJobLogParams{
		JobID: jobID,
	})
	if err != nil {
		return "", err
	}

	return response.Payload, nil
}

func ListWorkerPools() ([]*models.WorkerPool, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Jobservice.GetWorkerPools(ctx, &jobservice.GetWorkerPoolsParams{})
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func ListWorkers(poolID string) ([]*models.Worker, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Jobservice.GetWorkers(ctx, &jobservice.GetWorkersParams{
		PoolID: poolID,
	})
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}
