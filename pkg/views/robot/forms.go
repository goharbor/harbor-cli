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
package robot

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/huh"
)

// FormCreate asks for the basic fields of a new robot account.
func FormCreate(view *CreateView) error {
	duration := strconv.FormatInt(view.Duration, 10)
	if view.Duration == 0 {
		duration = "-1"
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Robot Name").
				Value(&view.Name).
				Validate(func(str string) error {
					if !isValidName(str) {
						return errors.New("invalid name: must start with a letter and contain only letters, digits, hyphen or underscore, no uppercase letters")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&view.Description),
			huh.NewInput().
				Title("Expiration").
				Value(&duration).
				Validate(func(str string) error {
					durationInt, err := strconv.Atoi(str)
					if err != nil {
						return errors.New("invalid expiration time: Enter expiration time in days")
					}
					if durationInt < -1 || durationInt == 0 {
						return errors.New("invalid expiration time: Enter -1 for no expiration or a positive integer for days")
					}
					view.Duration = int64(durationInt)
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}

// FormUpdate asks for the editable fields of an existing robot account.
func FormUpdate(view *UpdateView) error {
	duration := strconv.FormatInt(view.Duration, 10)

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Description").
				Value(&view.Description),
			huh.NewInput().
				Title("Expiration").
				Value(&duration).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("expiration cannot be empty")
					}
					dur, err := strconv.ParseInt(str, 10, 64)
					if err != nil {
						return errors.New("invalid expiration time: Enter expiration time in days")
					}
					view.Duration = dur
					return nil
				}),
			huh.NewConfirm().
				Title("Disable").
				Value(&view.Disable).
				Affirmative("yes").
				Negative("no"),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}

// FormSecret shows a robot's secret so the user can copy it.
func FormSecret(name, secret string) error {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Robot Name").
				Value(&name),
			huh.NewInput().
				Title("Robot Secret").
				Description("Copy the secret or press enter to copy to clipboard.").
				Value(&secret),
		),
	).WithTheme(huh.ThemeCharm()).Run()
}
