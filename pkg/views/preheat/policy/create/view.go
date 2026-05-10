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
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`

	// Provider related fields
	ProviderName string `json:"provider_name,omitempty"`

	// Filter related fields
	RepositoryFilter string `json:"repository_filter,omitempty"`
	TagFilter        string `json:"tag_filter,omitempty"`
	LabelFilter      string `json:"label_filter,omitempty"`

	// Trigger related fields
	TriggerType string `json:"trigger_type,omitempty"`
	CronString  string `json:"cron_string,omitempty"`
}

func CreatePreheatPolicyView(createView *CreateView, providers []*models.ProviderUnderProject) {
	if createView.TriggerType == "" {
		createView.TriggerType = "manual"
	}

	theme := huh.ThemeCharm()

	if len(providers) == 0 {
		log.Fatal("No P2P provider instances available for this project. Please create a provider instance first.")
	}

	providerOptions := make([]huh.Option[string], 0, len(providers))
	for _, p := range providers {
		if !p.Enabled {
			continue
		}
		label := fmt.Sprintf("%s (ID: %d)", p.Provider, p.ID)
		providerOptions = append(providerOptions, huh.NewOption(label, p.Provider))
	}

	if len(providerOptions) == 0 {
		log.Fatal("No enabled P2P provider instances available for this project.")
	}

	basicGroup := huh.NewGroup(
		huh.NewSelect[string]().
			Title("Provider").
			Options(providerOptions...).
			Value(&createView.ProviderName),
		huh.NewInput().
			Title("Policy Name").
			Value(&createView.Name).
			Validate(func(str string) error {
				if strings.TrimSpace(str) == "" {
					return errors.New("policy name cannot be empty")
				}
				return nil
			}),
		huh.NewInput().
			Title("Description").
			Value(&createView.Description),
		huh.NewConfirm().
			Title("Enabled").
			Value(&createView.Enabled).
			WithButtonAlignment(lipgloss.Left),
	)

	basicForm := huh.NewForm(basicGroup).WithTheme(theme)
	if err := basicForm.Run(); err != nil {
		log.Fatal(err)
	}

	filterGroup := huh.NewGroup(
		huh.NewInput().
			Title("Repositories").
			Description("Enter multiple comma separated repos, repo*, or **").
			Value(&createView.RepositoryFilter).
			Validate(func(str string) error {
				if strings.TrimSpace(str) == "" {
					return errors.New("repository filter cannot be empty")
				}
				return nil
			}),
		huh.NewInput().
			Title("Tags").
			Description("Enter multiple comma separated tags, tag*, or **").
			Value(&createView.TagFilter).
			Validate(func(str string) error {
				if strings.TrimSpace(str) == "" {
					return errors.New("tag filter cannot be empty")
				}
				return nil
			}),
		huh.NewInput().
			Title("Labels").
			Description("(optional)").
			Value(&createView.LabelFilter),
	).Title("Filters")

	triggerGroup := huh.NewGroup(
		huh.NewSelect[string]().
			Title("Trigger Type").
			Options(
				huh.NewOption("Manual", "manual"),
				huh.NewOption("Scheduled", "scheduled"),
				huh.NewOption("Event Based", "event_based"),
			).
			Value(&createView.TriggerType),
	)

	restForm := huh.NewForm(filterGroup, triggerGroup).WithTheme(theme)
	if err := restForm.Run(); err != nil {
		log.Fatal(err)
	}

	if createView.TriggerType == "scheduled" {
		schedulePreset := "none"
		presetForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Schedule").
					Description("Choose a schedule frequency for preheating").
					Options(
						huh.NewOption("None", "none"),
						huh.NewOption("Hourly", "hourly"),
						huh.NewOption("Daily", "daily"),
						huh.NewOption("Weekly", "weekly"),
						huh.NewOption("Custom", "custom"),
					).
					Value(&schedulePreset),
			),
		).WithTheme(theme)

		if err := presetForm.Run(); err != nil {
			log.Fatal(err)
		}

		if schedulePreset == "custom" {
			cronForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Cron String").
						Description("Schedule using 6-field cron format: seconds minutes hours day-month month day-week").
						Placeholder("0 0 0 * * *").
						Value(&createView.CronString).
						Validate(func(s string) error {
							if strings.TrimSpace(s) == "" {
								return errors.New("cron string cannot be empty for custom schedule")
							}

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
			createView.CronString = resolveSchedulePreset(schedulePreset, createView.CronString)
		} else {
			createView.CronString = resolveSchedulePreset(schedulePreset, "")
		}
	}
}

func resolveSchedulePreset(preset, cron string) string {
	switch preset {
	case "hourly":
		return "0 0 * * * *"
	case "daily":
		return "0 0 0 * * *"
	case "weekly":
		return "0 0 0 * * 0"
	case "custom":
		return strings.TrimSpace(cron)
	default:
		return ""
	}
}
