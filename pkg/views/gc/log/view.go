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

package log

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

// SelectGCJob displays an interactive list of GC jobs and returns the selected job ID
func SelectGCJob() (int64, error) {
	opts := api.ListFlags{Page: 1, PageSize: 100}
	history, err := api.GetGCHistory(opts)
	if err != nil {
		return 0, fmt.Errorf("failed to get GC history: %w", err)
	}

	if len(history) == 0 {
		return 0, fmt.Errorf("no GC jobs found")
	}

	options := buildOptions(history)

	var selectedID string
	theme := huh.ThemeCharm()

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a GC job to view logs").
				Description("Use arrow keys to navigate, press Enter to select").
				Options(options...).
				Value(&selectedID),
		),
	).WithTheme(theme).Run()

	if err != nil {
		return 0, fmt.Errorf("failed to select GC job: %w", err)
	}

	id, err := strconv.ParseInt(selectedID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid job ID: %w", err)
	}

	return id, nil
}

func buildOptions(history []*models.GCHistory) []huh.Option[string] {
	var options []huh.Option[string]
	for _, job := range history {
		creationTime, _ := utils.FormatCreatedTime(job.CreationTime.String())
		desc := fmt.Sprintf("ID: %d | Status: %s | Created: %s",
			job.ID, job.JobStatus, creationTime)
		options = append(options, huh.NewOption(desc, strconv.FormatInt(job.ID, 10)))
	}
	return options
}
