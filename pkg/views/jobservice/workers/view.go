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
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

// ListWorkers displays workers in a formatted table with pagination metadata.
func ListWorkers(workers []*models.Worker, page, pageSize, totalCount int64) {
	if len(workers) == 0 {
		fmt.Println("No workers found.")
		return
	}

	fmt.Printf("%-40s %-20s %-20s %-40s %-30s %-30s\n",
		"ID", "POOL_ID", "JOB_NAME", "JOB_ID", "START_AT", "CHECKIN_AT")
	fmt.Printf("%-40s %-20s %-20s %-40s %-30s %-30s\n",
		"--", "-------", "--------", "------", "--------", "---------")

	busyCount := 0
	for _, worker := range workers {
		id := worker.ID
		if id == "" {
			id = "-"
		}

		poolID := worker.PoolID
		if poolID == "" {
			poolID = "-"
		}

		jobName := worker.JobName
		if jobName == "" {
			jobName = "-"
		}

		jobID := worker.JobID
		if jobID == "" {
			jobID = "-"
		} else {
			busyCount++
		}

		startAt := "-"
		if worker.StartAt != nil {
			startAt = fmt.Sprintf("%v", worker.StartAt)
		}

		checkinAt := "-"
		if worker.CheckinAt != nil {
			checkinAt = fmt.Sprintf("%v", worker.CheckinAt)
		}

		fmt.Printf("%-40s %-20s %-20s %-40s %-30s %-30s\n",
			id, poolID, jobName, jobID, startAt, checkinAt)
	}

	fmt.Printf("\nPage: %d  Page Size: %d  Returned: %d  Total: %d  Busy: %d\n", page, pageSize, len(workers), totalCount, busyCount)
}
