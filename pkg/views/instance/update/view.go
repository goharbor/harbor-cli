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
package update

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func UpdateInstanceView(instance *models.Instance) error {
	theme := huh.ThemeCharm()

	authUsername := ""
	authPassword := ""
	authToken := ""

	if instance.AuthInfo != nil {
		switch strings.ToUpper(instance.AuthMode) {
		case "BASIC":
			authUsername = instance.AuthInfo["username"]
			authPassword = instance.AuthInfo["password"]
		case "OAUTH":
			authToken = instance.AuthInfo["token"]
		}
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Provider (cannot be changed)").
				Description(instance.Vendor),
			huh.NewInput().
				Title("Name").
				Value(&instance.Name).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("name cannot be empty or only spaces")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&instance.Description),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Endpoint").
				Value(&instance.Endpoint).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("endpoint cannot be empty or only spaces")
					}
					formattedURL := utils.FormatUrl(str)
					if err := utils.ValidateURL(formattedURL); err != nil {
						return err
					}
					instance.Endpoint = formattedURL
					return nil
				}),
			huh.NewConfirm().
				Title("Enable").
				Value(&instance.Enabled).
				Affirmative("yes").
				Negative("no").
				WithButtonAlignment(lipgloss.Left),
			huh.NewConfirm().
				Title("Skip Certificate Verification").
				Value(&instance.Insecure).
				Affirmative("yes").
				Negative("no").
				WithButtonAlignment(lipgloss.Left),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Auth Mode").
				Options(
					huh.NewOption("None", "NONE"),
					huh.NewOption("Basic", "BASIC"),
					huh.NewOption("OAuth", "OAUTH"),
				).
				Value(&instance.AuthMode),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Value(&authUsername).
				Validate(func(str string) error {
					if instance.AuthMode != "BASIC" {
						return nil
					}
					if strings.TrimSpace(str) == "" {
						return errors.New("username cannot be empty or only spaces")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&authPassword).
				Validate(func(str string) error {
					if instance.AuthMode != "BASIC" {
						return nil
					}
					if strings.TrimSpace(str) == "" {
						return errors.New("password cannot be empty or only spaces")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return instance.AuthMode == "NONE" || instance.AuthMode == "OAUTH"
		}),
		huh.NewGroup(
			huh.NewInput().
				Title("Token").
				EchoMode(huh.EchoModePassword).
				Value(&authToken).
				Validate(func(str string) error {
					if instance.AuthMode != "OAUTH" {
						return nil
					}
					if strings.TrimSpace(str) == "" {
						return errors.New("token cannot be empty or only spaces")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return instance.AuthMode == "NONE" || instance.AuthMode == "BASIC"
		}),
	).WithTheme(theme).Run()
	if err != nil {
		return err
	}

	instance.AuthMode = strings.ToUpper(strings.TrimSpace(instance.AuthMode))

	switch instance.AuthMode {
	case "NONE":
		instance.AuthInfo = nil
	case "BASIC":
		instance.AuthInfo = map[string]string{
			"username": strings.TrimSpace(authUsername),
			"password": authPassword,
		}
	case "OAUTH":
		instance.AuthInfo = map[string]string{
			"token": strings.TrimSpace(authToken),
		}
	}

	return nil
}
