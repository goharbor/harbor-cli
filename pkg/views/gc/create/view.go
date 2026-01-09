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
package gccreate

import (
	"fmt"
	"regexp"

	"github.com/charmbracelet/huh"
	"github.com/go-openapi/strfmt"
)

type CreateView struct {
	Type              string
	Cron              string
	NextScheduledTime strfmt.DateTime
	Parameters        map[string]any
}

func CreateScheduleView(createView *CreateView) error {
	theme := huh.ThemeCharm()

	typeOpts := []huh.Option[string]{
		huh.NewOption("Hourly", "Hourly"),
		huh.NewOption("Daily", "Daily"),
		huh.NewOption("Weekly", "Weekly"),
		huh.NewOption("None", "None"),
		huh.NewOption("Custom", "Custom"),
	}
	var nextScheduleTime string
	var addParam bool

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Schedule Type").
				Value(&createView.Type).
				Options(typeOpts...),
			huh.NewInput().
				Title("Next Scheduled Time").
				Value(&nextScheduleTime).
				Description("Example: 2026-01-04T10:15:30Z").
				Validate(func(s string) error {
					_, err := strfmt.ParseDateTime(s)
					if err != nil {
						return err
					}

					return nil
				}),
		),
	).WithTheme(theme).Run()
	if err != nil {
		return err
	}

	if createView.Type == "Custom" {
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Cron").
					Value(&createView.Cron).
					Description("Example: 0 9 * * 1-5").
					Validate(func(s string) error {
						cronBasic := regexp.MustCompile(`^(\S+\s+){4}\S+$`)
						if cronBasic.MatchString(s) {
							return fmt.Errorf("invalid cron format")
						}

						return nil
					}),
			),
		).WithTheme(theme).Run()
		if err != nil {
			return err
		}
	}

	for !addParam {
		var (
			k string
			v string
		)

		err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Key").
					Value(&k),
				huh.NewInput().
					Title("Value").
					Value(&v),
			),
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add Another Parameter?").
					Value(&addParam),
			),
		).WithTheme(theme).Run()
		if err != nil {
			return err
		}

		createView.Parameters[k] = v
	}

	dt, err := strfmt.ParseDateTime(nextScheduleTime)
	if err != nil {
		return err
	}

	createView.NextScheduledTime = dt

	return nil
}
