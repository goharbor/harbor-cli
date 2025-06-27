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
	"strconv"
	"unicode"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	Description string             `json:"description,omitempty"`
	Disable     bool               `json:"disable,omitempty"`
	Duration    int64              `json:"duration,omitempty"`
	Level       string             `json:"level,omitempty"`
	Name        string             `json:"name,omitempty"`
	Permissions []*RobotPermission `json:"permissions"`
	Secret      string             `json:"secret,omitempty"`
	ProjectName string
}

type RobotPermission struct {
	Access    []*models.Access `json:"access"`
	Kind      string           `json:"kind,omitempty"`
	Namespace string           `json:"namespace,omitempty"`
}

type Access struct {
	Action   string `json:"action,omitempty"`
	Effect   string `json:"effect,omitempty"`
	Resource string `json:"resource,omitempty"`
}

func CreateRobotView(createView *CreateView) {
	duration := strconv.FormatInt(createView.Duration, 10)
	if createView.Duration == 0 {
		duration = "-1"
	}

	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Robot Name").
				Value(&createView.Name).
				Validate(func(str string) error {
					if !isValidName(str) {
						return errors.New("invalid name: must start with a letter and contain only letters, digits, hyphen or underscore, no uppercase letters")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description),
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
					dur, err := strconv.ParseInt(str, 10, 64)
					if err != nil {
						return errors.New("invalid expiration time: Enter expiration time in days")
					}
					createView.Duration = dur
					return nil
				}),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func CreateRobotSecretView(name string, secret string) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Robot Name").
				Value(&name),
			huh.NewInput().
				Title("Robot Secret").
				Description("Copy the secret or press enter to copy to clipboard.").
				Value(&secret),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func isValidName(s string) bool {
	if s == "" {
		return false
	}
	if len(s) > 0 && !unicode.IsLetter(rune(s[0])) {
		return false
	}
	for _, r := range s {
		if unicode.IsUpper(r) {
			return false
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false
		}
	}
	return true
}
