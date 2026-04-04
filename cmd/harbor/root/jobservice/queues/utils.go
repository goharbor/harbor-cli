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
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/api"
)

func shouldIncludeQueueForAction(action string, paused bool) bool {
	switch strings.ToLower(action) {
	case "resume":
		return paused
	case "pause":
		return !paused
	default:
		return true
	}
}

func executeQueueAction(action string, jobTypes []string) error {
	normalizedTypes := normalizeJobTypes(jobTypes)
	if len(normalizedTypes) == 0 {
		return fmt.Errorf("at least one job type must be provided")
	}

	for _, jobType := range normalizedTypes {
		fmt.Printf("%s queue type '%s'...\n", actionLabel(action), jobType)
		err := api.ActionJobQueue(strings.ToUpper(jobType), action)
		if err != nil {
			return fmt.Errorf("failed to %s queue '%s': %w", action, jobType, err)
		}
		fmt.Printf("✓ Queue '%s' %sd successfully.\n", jobType, action)
	}

	return nil
}

func normalizeJobTypes(jobTypes []string) []string {
	cleanedTypes := make([]string, 0, len(jobTypes))
	seen := make(map[string]struct{}, len(jobTypes))

	for _, rawType := range jobTypes {
		for _, splitType := range strings.Split(rawType, ",") {
			trimmedType := strings.TrimSpace(splitType)
			if trimmedType == "" {
				continue
			}

			if strings.EqualFold(trimmedType, "all") {
				return []string{"all"}
			}

			key := strings.ToLower(trimmedType)
			if _, exists := seen[key]; exists {
				continue
			}

			seen[key] = struct{}{}
			cleanedTypes = append(cleanedTypes, trimmedType)
		}
	}

	return cleanedTypes
}

func actionLabel(action string) string {
	if action == "" {
		return "Updating"
	}

	lower := strings.ToLower(action)
	return strings.ToUpper(lower[:1]) + lower[1:]
}

// selectQueueTypes shows an interactive multi-selector for queue types
func selectQueueTypes(action string) ([]string, error) {
	response, err := api.ListJobQueues()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve job queues: %w", err)
	}

	if response == nil || response.Payload == nil || len(response.Payload) == 0 {
		return nil, fmt.Errorf("no job queues available")
	}

	filteredQueues := make([]*struct {
		JobType string
		Count   int64
	}, 0, len(response.Payload))

	for _, queue := range response.Payload {
		if queue == nil {
			continue
		}
		if shouldIncludeQueueForAction(action, queue.Paused) {
			filteredQueues = append(filteredQueues, &struct {
				JobType string
				Count   int64
			}{
				JobType: queue.JobType,
				Count:   queue.Count,
			})
		}
	}

	if len(filteredQueues) == 0 {
		switch action {
		case "resume":
			return nil, fmt.Errorf("no paused queues available to resume")
		case "pause":
			return nil, fmt.Errorf("all queues are already paused")
		default:
			return nil, fmt.Errorf("no job queues available to %s", action)
		}
	}

	options := make([]huh.Option[string], len(filteredQueues)+1)
	options[0] = huh.NewOption("all", "all")

	for i, queue := range filteredQueues {
		label := fmt.Sprintf("%s (pending: %d)", queue.JobType, queue.Count)
		options[i+1] = huh.NewOption(label, queue.JobType)
	}

	var selected []string
	theme := huh.ThemeCharm()
	keymap := huh.NewDefaultKeyMap()
	keymap.Quit = key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	)

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(fmt.Sprintf("Select queue type(s) to %s (press q to cancel)", action)).
				Options(options...).
				Value(&selected),
		),
	).WithTheme(theme).WithKeyMap(keymap).Run()

	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, errors.New("operation cancelled")
		}
		return nil, err
	}

	selected = normalizeJobTypes(selected)
	if len(selected) == 0 {
		return nil, fmt.Errorf("at least one queue type must be selected")
	}

	return selected, nil
}
