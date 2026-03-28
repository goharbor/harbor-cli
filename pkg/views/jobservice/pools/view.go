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
package pools

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

// ListPools displays worker pools in a formatted table.
func ListPools(items []*models.WorkerPool) {
	if len(items) == 0 {
		fmt.Println("No worker pools found.")
		return
	}

	fmt.Printf("%-20s %-8s %-30s %-30s %-12s %-30s\n", "POOL_ID", "PID", "START_AT", "HEARTBEAT_AT", "CONCURRENCY", "HOST")
	fmt.Printf("%-20s %-8s %-30s %-30s %-12s %-30s\n", "-------", "---", "--------", "------------", "-----------", "----")

	for _, pool := range items {
		if pool == nil {
			continue
		}

		startAt := fmt.Sprintf("%v", pool.StartAt)
		heartbeatAt := fmt.Sprintf("%v", pool.HeartbeatAt)

		fmt.Printf("%-20s %-8d %-30s %-30s %-12d %-30s\n",
			pool.WorkerPoolID,
			pool.Pid,
			startAt,
			heartbeatAt,
			pool.Concurrency,
			pool.Host,
		)
	}

	fmt.Printf("\nTotal: %d worker pool(s)\n", len(items))
}
