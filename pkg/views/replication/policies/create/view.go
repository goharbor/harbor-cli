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
package create

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	Enabled           bool   `json:"enabled,omitempty"`
	ReplicationMode   string `json:"mode,omitempty"`
	Override          bool   `json:"override,omitempty"`
	ReplicateDeletion bool   `json:"replicate_deletion,omitempty"`
	CopyByChunk       bool   `json:"copy_by_chunk,omitempty"`
	Speed             string `json:"speed,omitempty"`

	// Trigger related fields
	TriggerType           string `json:"trigger_type,omitempty"`
	DeleteRemoteResources bool   `json:"delete_remote_resources,omitempty"`
	CronString            string `json:"cron_string,omitempty"`
}

func CreateRPolicyView(createView *CreateView) {
	if createView.ReplicationMode == "" {
		createView.ReplicationMode = "Pull"
	}
	if createView.TriggerType == "" {
		createView.TriggerType = "manual"
	}

	createView.Override = true
	theme := huh.ThemeCharm()
	basicGroup := huh.NewGroup(
		huh.NewInput().
			Title("Replication Policy Name").
			Value(&createView.Name).
			Validate(func(str string) error {
				if str == "" {
					return errors.New("name cannot be empty")
				}
				return nil
			}),
		huh.NewInput().
			Title("Description").
			Value(&createView.Description),
		huh.NewSelect[string]().
			Title("Replication Mode").
			Description("Choose whether to pull from or push to an external registry").
			Options(
				huh.NewOption("Pull (External → Harbor)", "Pull"),
				huh.NewOption("Push (Harbor → External)", "Push"),
			).
			Value(&createView.ReplicationMode),
	)

	triggerGroup := huh.NewGroup(
		huh.NewSelect[string]().
			Title("Trigger Mode").
			Description("When should replication occur?").
			Options(
				huh.NewOption("Manual", "manual"),
				huh.NewOption("Event Based", "event_based"),
				huh.NewOption("Scheduled", "scheduled"),
			).
			Value(&createView.TriggerType),
	)

	advancedGroup := huh.NewGroup(
		huh.NewConfirm().
			Title("Override").
			Description("Replace artifacts on destination if they already exist").
			Value(&createView.Override).
			WithButtonAlignment(lipgloss.Left),
		huh.NewConfirm().
			Title("Replicate Deletion").
			Description("Synchronize deletion operations between registries").
			Value(&createView.ReplicateDeletion).
			WithButtonAlignment(lipgloss.Left),
		huh.NewConfirm().
			Title("Copy By Chunk").
			Description("Transfer artifacts in smaller chunks for better reliability").
			Value(&createView.CopyByChunk).
			WithButtonAlignment(lipgloss.Left),
		huh.NewInput().
			Title("Speed Limit").
			Description("Maximum speed in KB/s (-1 = unlimited)").
			Placeholder("-1").
			Value(&createView.Speed).
			Validate(func(s string) error {
				if s == "" {
					return nil
				}
				speed, err := strconv.ParseInt(s, 10, 32)
				if err != nil {
					return fmt.Errorf("speed must be a valid number")
				}
				if speed < -1 {
					return fmt.Errorf("speed cannot be negative")
				}
				return nil
			}),
		huh.NewConfirm().
			Title("Enabled").
			Description("Activate replication policy after creation").
			Value(&createView.Enabled).
			WithButtonAlignment(lipgloss.Left),
	)

	form := huh.NewForm(basicGroup, triggerGroup, advancedGroup)
	form.WithTheme(theme)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	if createView.TriggerType == "event_based" {
		eventForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Delete remote resources when locally deleted").
					Description("When artifacts are deleted locally, also delete them on the remote registry").
					Value(&createView.DeleteRemoteResources).
					WithButtonAlignment(lipgloss.Left),
			),
		).WithTheme(theme)

		if err := eventForm.Run(); err != nil {
			log.Fatal(err)
		}
	} else if createView.TriggerType == "scheduled" {
		cronForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Cron String").
					Description("Schedule using 6-field cron format: seconds minutes hours day-month month day-week").
					Placeholder("0 0 0 * * *"). // At midnight (00:00:00) every day
					Value(&createView.CronString).
					Validate(func(s string) error {
						if s == "" {
							return errors.New("cron string cannot be empty for scheduled trigger")
						}

						// Basic validation for 6-field cron format
						fields := strings.Fields(s)
						if len(fields) != 6 {
							return fmt.Errorf("cron must have exactly 6 fields (found %d): seconds minutes hours day-month month day-week", len(fields))
						}

						return nil
					}),
			),
		).WithTheme(theme)

		if err := cronForm.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
