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

	fmt.Printf("%-36s %-20s %-10s %-20s %-20s\n", "POOL_ID", "HOST", "PID", "CONCURRENCY", "START_AT")
	fmt.Printf("%-36s %-20s %-10s %-20s %-20s\n", "-------", "----", "---", "-----------", "--------")

	for _, p := range items {
		if p == nil {
			continue
		}
		fmt.Printf("%-36s %-20s %-10d %-20d %-20s\n", p.WorkerPoolID, p.Host, p.Pid, p.Concurrency, p.StartAt.String())
	}

	fmt.Printf("\nTotal: %d pool(s)\n", len(items))
}
