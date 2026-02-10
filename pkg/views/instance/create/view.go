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
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	Vendor      string
	Name        string
	Description string
	Endpoint    string
	AuthMode    string
	AuthInfo    map[string]string
	Enabled     bool
	Insecure    bool
}

func CreateInstanceView(createView *CreateView) {
	cv := CreateView{
		AuthInfo: map[string]string{
			"username": "",
			"password": "",
			"token":    "",
		},
	}
	username := cv.AuthInfo["username"]
	password := cv.AuthInfo["password"]
	token := cv.AuthInfo["token"]
	theme := huh.ThemeCharm()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Provider").
				Options(
					huh.NewOption("Dragonfly", "dragonfly"),
					huh.NewOption("Kraken", "kraken"),
				).
				Value(&createView.Vendor),
			huh.NewInput().
				Title("Name").
				Value(&createView.Name).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("name cannot be empty or only spaces")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Endpoint").
				Value(&createView.Endpoint).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("endpoint cannot be empty or only spaces")
					}
					formattedUrl := utils.FormatUrl(str)
					if err := utils.ValidateURL(formattedUrl); err != nil {
						return err
					}
					createView.Endpoint = formattedUrl
					return nil
				}),
			huh.NewConfirm().
				Title("Enable").
				Value(&createView.Enabled).
				Affirmative("yes").
				Negative("no"),
			huh.NewConfirm().
				Title("Verify Cert").
				Value(&createView.Insecure).
				Affirmative("yes").
				Negative("no"),
		),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Auth Mode").
				Options(
					huh.NewOption("None", "NONE"),
					huh.NewOption("Basic", "BASIC"),
					huh.NewOption("OAuth", "OAUTH"),
				).
				Value(&createView.AuthMode),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Value(&username).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("username cannot be empty or only spaces")
					}
					if isValid := utils.ValidateUserName(str); !isValid {
						return errors.New("please enter correct username format")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&password).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("password cannot be empty or only spaces")
					}
					if err := utils.ValidatePassword(str); err != nil {
						return err
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return createView.AuthMode == "NONE" || createView.AuthMode == "OAUTH"
		}),
		huh.NewGroup(
			huh.NewInput().
				Title("Token").
				Value(&token).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("token cannot be empty or only spaces")
					}
					return nil
				}),
		).WithHideFunc(func() bool {
			return createView.AuthMode == "NONE" || createView.AuthMode == "BASIC"
		}),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
