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

// ListWorkers displays workers in a formatted table.
func ListWorkers(items []*models.Worker) {
	if len(items) == 0 {
		fmt.Println("No workers found.")
		return
	}

	fmt.Printf("%-36s %-36s %-15s %-20s\n", "WORKER_ID", "POOL_ID", "JOB_ID", "JOB_NAME")
	fmt.Printf("%-36s %-36s %-15s %-20s\n", "---------", "-------", "------", "--------")

	for _, w := range items {
		if w == nil {
			continue
		}
		fmt.Printf("%-36s %-36s %-15s %-20s\n", w.ID, w.PoolID, w.JobID, w.JobName)
	}

	fmt.Printf("\nTotal: %d worker(s)\n", len(items))
}
