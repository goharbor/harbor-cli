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
package jobservice

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/base/loadingtable"
	"github.com/goharbor/harbor-cli/pkg/views/base/tablelist"
)

var (
	queueColumns = []table.Column{
		{Title: "Job Type", Width: tablelist.WidthXXL},
		{Title: "Pending Jobs", Width: tablelist.WidthL},
		{Title: "Latency (s)", Width: tablelist.WidthL},
		{Title: "Paused", Width: tablelist.WidthS},
	}
	poolColumns = []table.Column{
		{Title: "Worker Pool ID", Width: tablelist.WidthXXL},
		{Title: "Pid", Width: tablelist.WidthS},
		{Title: "Started At", Width: tablelist.WidthL},
		{Title: "Heartbeat At", Width: tablelist.WidthL},
		{Title: "Concurrency", Width: tablelist.WidthL},
	}
	workerColumns = []table.Column{
		{Title: "Worker ID", Width: tablelist.WidthXXL},
		{Title: "Pool ID", Width: tablelist.WidthXXL},
		{Title: "Job Name", Width: tablelist.WidthL},
		{Title: "Job ID", Width: tablelist.WidthXXL},
		{Title: "Started At", Width: tablelist.WidthL},
		{Title: "Checked In At", Width: tablelist.WidthL},
	}
)

var ErrUserAborted = errors.New("user aborted selection")

func jobQueueRows(queues []*models.JobQueue) []table.Row {
	var rows []table.Row
	for _, queue := range queues {
		paused := "No"
		if queue.Paused {
			paused = "Yes"
		}
		rows = append(rows, table.Row{
			queue.JobType,
			strconv.FormatInt(queue.Count, 10),
			strconv.FormatInt(queue.Latency, 10),
			paused,
		})
	}
	return rows
}

func workerPoolRows(pools []*models.WorkerPool) []table.Row {
	var rows []table.Row
	for _, pool := range pools {
		startedAt, _ := utils.FormatCreatedTime(pool.StartAt.String())
		heartbeatAt, _ := utils.FormatCreatedTime(pool.HeartbeatAt.String())
		rows = append(rows, table.Row{
			pool.WorkerPoolID,
			strconv.FormatInt(pool.Pid, 10),
			startedAt,
			heartbeatAt,
			strconv.FormatInt(pool.Concurrency, 10),
		})
	}
	return rows
}

func workerRows(workers []*models.Worker) []table.Row {
	var rows []table.Row
	for _, worker := range workers {
		startedAt := ""
		if worker.StartAt != nil {
			startedAt, _ = utils.FormatCreatedTime(worker.StartAt.String())
		}
		checkedInAt := ""
		if worker.CheckinAt != nil {
			checkedInAt, _ = utils.FormatCreatedTime(worker.CheckinAt.String())
		}
		rows = append(rows, table.Row{
			worker.ID,
			worker.PoolID,
			worker.JobName,
			worker.JobID,
			startedAt,
			checkedInAt,
		})
	}
	return rows
}

func ListJobQueuesAsync() error {
	fetcher := func() tea.Msg {
		queues, err := api.ListJobQueues()
		if err != nil {
			return loadingtable.FetchMsg{Err: err}
		}
		return loadingtable.FetchMsg{Rows: jobQueueRows(queues)}
	}

	m := loadingtable.NewModel("Job Queues", queueColumns, fetcher)

	_, err := tea.NewProgram(m).Run()
	return err
}

func ListWorkerPoolsAsync() error {
	fetcher := func() tea.Msg {
		pools, err := api.ListWorkerPools()
		if err != nil {
			return loadingtable.FetchMsg{Err: err}
		}
		return loadingtable.FetchMsg{Rows: workerPoolRows(pools)}
	}

	m := loadingtable.NewModel("Worker Pools", poolColumns, fetcher)

	_, err := tea.NewProgram(m).Run()
	return err
}

func ListWorkersAsync(poolID string) error {
	title := "Workers"
	if poolID != "" {
		title = fmt.Sprintf("Workers in Pool: %s", poolID)
	}

	fetcher := func() tea.Msg {
		workers, err := api.ListWorkers(poolID)
		if err != nil {
			return loadingtable.FetchMsg{Err: err}
		}
		return loadingtable.FetchMsg{Rows: workerRows(workers)}
	}

	m := loadingtable.NewModel(title, workerColumns, fetcher)

	_, err := tea.NewProgram(m).Run()
	return err
}

func SelectQueueAsync(title string) (string, error) {
	fetcher := func() tea.Msg {
		queues, err := api.ListJobQueues()
		if err != nil {
			return loadingtable.FetchMsg{Err: err}
		}
		rows := []table.Row{{"all", "", "", ""}}
		rows = append(rows, jobQueueRows(queues)...)
		return loadingtable.FetchMsg{Rows: rows}
	}

	m := loadingtable.NewModel(title, queueColumns, fetcher)

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return "", err
	}

	if model, ok := p.(loadingtable.Model); ok {
		if model.Aborted {
			return "", ErrUserAborted
		}
		if !model.Selected {
			return "", fmt.Errorf("no queue selected")
		}
		return model.Choice[0], nil
	}

	return "", fmt.Errorf("unexpected selection result")
}

func SelectPoolAsync(title string) (string, error) {
	fetcher := func() tea.Msg {
		pools, err := api.ListWorkerPools()
		if err != nil {
			return loadingtable.FetchMsg{Err: err}
		}
		rows := []table.Row{{"all", "", "", "", ""}}
		rows = append(rows, workerPoolRows(pools)...)
		return loadingtable.FetchMsg{Rows: rows}
	}

	m := loadingtable.NewModel(title, poolColumns, fetcher)

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return "", err
	}

	if model, ok := p.(loadingtable.Model); ok {
		if model.Aborted {
			return "", ErrUserAborted
		}
		if !model.Selected {
			return "", fmt.Errorf("no pool selected")
		}
		return model.Choice[0], nil
	}

	return "", fmt.Errorf("unexpected selection result")
}

func SelectRunningJobAsync() (string, error) {
	fetcher := func() tea.Msg {
		workers, err := api.ListWorkers("")
		if err != nil {
			return loadingtable.FetchMsg{Err: err}
		}
		var rows []table.Row
		for _, w := range workers {
			if w.JobID != "" {
				rows = append(rows, table.Row{w.JobName, w.JobID, w.ID})
			}
		}
		if len(rows) == 0 {
			return loadingtable.FetchMsg{Err: errors.New("no running jobs found")}
		}
		return loadingtable.FetchMsg{Rows: rows}
	}

	columns := []table.Column{
		{Title: "Job Name", Width: tablelist.WidthL},
		{Title: "Job ID", Width: tablelist.WidthXXL},
		{Title: "Worker ID", Width: tablelist.WidthXXL},
	}

	m := loadingtable.NewModel("Select a Running Job", columns, fetcher)

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return "", err
	}

	if model, ok := p.(loadingtable.Model); ok {
		if model.Aborted {
			return "", ErrUserAborted
		}
		if !model.Selected {
			return "", fmt.Errorf("no job selected")
		}
		return model.Choice[1], nil
	}

	return "", fmt.Errorf("unexpected selection result")
}
