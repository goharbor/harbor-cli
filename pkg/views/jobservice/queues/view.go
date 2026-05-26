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
package queues

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

// ListQueues displays job queues in a formatted table.
func ListQueues(items []*models.JobQueue) {
	if len(items) == 0 {
		fmt.Println("No job queues found.")
		return
	}

	fmt.Printf("%-25s %-12s %-12s %-10s\n", "JOB_TYPE", "COUNT", "LATENCY(s)", "PAUSED")
	fmt.Printf("%-25s %-12s %-12s %-10s\n", "--------", "-----", "----------", "------")

	for _, queue := range items {
		if queue == nil {
			continue
		}
		fmt.Printf("%-25s %-12d %-12d %-10t\n", queue.JobType, queue.Count, queue.Latency, queue.Paused)
	}

	fmt.Printf("\nTotal: %d queue(s)\n", len(items))
}
