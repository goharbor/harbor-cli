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
package schedules

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

// ListSchedules displays schedule tasks with pagination metadata.
func ListSchedules(items []*models.ScheduleTask, page, pageSize, totalCount int64) {
	if len(items) == 0 {
		fmt.Println("No schedules found.")
		return
	}

	fmt.Printf("%-8s %-18s %-22s %-30s %-20s\n", "ID", "VENDOR_TYPE", "VENDOR_ID", "CRON", "UPDATE_TIME")
	fmt.Printf("%-8s %-18s %-22s %-30s %-20s\n", "--", "-----------", "---------", "----", "-----------")

	for _, task := range items {
		if task == nil {
			continue
		}
		fmt.Printf("%-8d %-18s %-22d %-30s %-20s\n", task.ID, task.VendorType, task.VendorID, task.Cron, task.UpdateTime.String())
	}

	fmt.Printf("\nPage: %d  Page Size: %d  Returned: %d  Total: %d\n", page, pageSize, len(items), totalCount)
}

// PrintScheduleStatus displays the scheduler paused/running state.
func PrintScheduleStatus(status *models.SchedulerStatus) {
	if status == nil {
		fmt.Println("Scheduler status: unknown")
		return
	}

	if status.Paused {
		fmt.Println("Scheduler status: paused")
		return
	}

	fmt.Println("Scheduler status: running")
}
