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
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/jobservice/queues"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// QueuesCommand creates the queues subcommand
func QueuesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queues",
		Short: "Manage job queues (list, stop, pause, resume)",
		Long:  "List job queues and perform actions on them (stop/pause/resume).",
	}

	cmd.AddCommand(ListCommand(), StopCommand(), PauseCommand(), ResumeCommand())

	return cmd
}

// ListCommand lists all job queues
func ListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all job queues",
		Long:    "Display all job queues with their pending job counts and latency.",
		Example: "harbor jobservice queues list",
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := api.ListJobQueues()
			if err != nil {
				return fmt.Errorf("failed to retrieve job queues: %w", err)
			}

			if response == nil || response.Payload == nil || len(response.Payload) == 0 {
				fmt.Println("No job queues found.")
				return nil
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				return utils.PrintFormat(response.Payload, formatFlag)
			}

			queues.ListQueues(response.Payload)
			return nil
		},
	}

	return cmd
}

// StopCommand stops a job queue
func StopCommand() *cobra.Command {
	var jobTypes []string
	var interactive bool

	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stop queue(s) (--type or --interactive)",
		Long:    "Stop a job queue or all queues.",
		Example: "harbor jobservice queues stop --type REPLICATION\nharbor jobservice queues stop --type REPLICATION --type RETENTION\nharbor jobservice queues stop --type all",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(jobTypes) == 0 && !interactive {
				interactive = true
			}

			if interactive {
				selectedTypes, err := selectQueueTypes("stop")
				if err != nil {
					return err
				}
				jobTypes = selectedTypes
			}

			if len(jobTypes) == 0 {
				return fmt.Errorf("at least one job type must be specified with --type or interactive mode")
			}

			return executeQueueAction("stop", jobTypes)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVar(&jobTypes, "type", nil, "Job type(s) to stop (repeat flag or comma-separate values; use 'all' for all queues)")
	flags.BoolVarP(&interactive, "interactive", "i", false, "Interactive mode to choose queue type(s) instead of passing --type")

	return cmd
}

// PauseCommand pauses a job queue
func PauseCommand() *cobra.Command {
	var jobTypes []string
	var interactive bool

	cmd := &cobra.Command{
		Use:     "pause",
		Short:   "Pause queue(s) (--type or --interactive)",
		Long:    "Pause a job queue or all queues.",
		Example: "harbor jobservice queues pause --type REPLICATION\nharbor jobservice queues pause --type REPLICATION --type RETENTION\nharbor jobservice queues pause --type all",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(jobTypes) == 0 && !interactive {
				interactive = true
			}

			if interactive {
				selectedTypes, err := selectQueueTypes("pause")
				if err != nil {
					return err
				}
				jobTypes = selectedTypes
			}

			if len(jobTypes) == 0 {
				return fmt.Errorf("at least one job type must be specified with --type or interactive mode")
			}

			return executeQueueAction("pause", jobTypes)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVar(&jobTypes, "type", nil, "Job type(s) to pause (repeat flag or comma-separate values; use 'all' for all queues)")
	flags.BoolVarP(&interactive, "interactive", "i", false, "Interactive mode to choose queue type(s) instead of passing --type")

	return cmd
}

// ResumeCommand resumes a job queue
func ResumeCommand() *cobra.Command {
	var jobTypes []string
	var interactive bool

	cmd := &cobra.Command{
		Use:     "resume",
		Short:   "Resume queue(s) (--type or --interactive)",
		Long:    "Resume a paused job queue or all queues.",
		Example: "harbor jobservice queues resume --type REPLICATION\nharbor jobservice queues resume --type REPLICATION --type RETENTION\nharbor jobservice queues resume --type all",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(jobTypes) == 0 && !interactive {
				interactive = true
			}

			if interactive {
				selectedTypes, err := selectQueueTypes("resume")
				if err != nil {
					return err
				}
				jobTypes = selectedTypes
			}

			if len(jobTypes) == 0 {
				return fmt.Errorf("at least one job type must be specified with --type or interactive mode")
			}

			return executeQueueAction("resume", jobTypes)
		},
	}

	flags := cmd.Flags()
	flags.StringSliceVar(&jobTypes, "type", nil, "Job type(s) to resume (repeat flag or comma-separate values; use 'all' for all queues)")
	flags.BoolVarP(&interactive, "interactive", "i", false, "Interactive mode to choose queue type(s) instead of passing --type")

	return cmd
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
