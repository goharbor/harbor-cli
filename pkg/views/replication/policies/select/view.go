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
package replication

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/base/selection"
)

func ReplicationPoliciesList(policies []*models.ReplicationPolicy, choice chan<- int64) {
	policyNameIDsMap := make(map[string]int64, len(policies))
	for _, p := range policies {
		policyNameIDsMap[p.Name] = p.ID
	}

	itemsList := make([]list.Item, len(policies))
	for i, p := range policies {
		itemsList[i] = selection.Item(p.Name)
	}

	m := selection.NewModel(itemsList, "Replication Policy")

	p, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if p, ok := p.(selection.Model); ok {
		choice <- policyNameIDsMap[p.Choice]
	}
}
