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
package login

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type LoginView struct {
	Server   string
	Username string
	Password string
	Name     string
	Config   string
}

func CreateView(loginView *LoginView) {
	theme := huh.ThemeCharm()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server").
				Description("Server address eg. demo.goharbor.io").
				Value(&loginView.Server).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("server cannot be empty or only spaces")
					}
					formattedUrl := utils.FormatUrl(str)
					if _, err := url.ParseRequestURI(formattedUrl); err != nil {
						return errors.New("please enter the correct server format")
					}
					return nil
				}),
			huh.NewInput().
				Title("User Name").
				Value(&loginView.Username).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("username cannot be empty or only spaces")
					}
					if isValid := utils.ValidateName(str); !isValid {
						return errors.New("please enter correct username format")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&loginView.Password).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("password cannot be empty or only spaces")
					}
					if err := utils.ValidatePassword(str); err != nil {
						return err
					}
					return nil
				}),
			huh.NewInput().
				Title("Name of Credential").
				Value(&loginView.Name).
				Description("Name of credential to be stored in the harbor config file.").
				PlaceholderFunc(func() string {
					return fmt.Sprintf("%s@%s", loginView.Username, utils.SanitizeServerAddress(loginView.Server))
				}, &loginView).
				SuggestionsFunc(func() []string {
					return []string{
						fmt.Sprintf("%s@%s", loginView.Username, utils.SanitizeServerAddress(loginView.Server)),
					}
				}, &loginView).
				Validate(func(str string) error {
					if str == "" {
						loginView.Name = fmt.Sprintf("%s@%s", loginView.Username, utils.SanitizeServerAddress(loginView.Server))
						return nil
					}
					return nil
				}),
		),
	).WithTheme(theme).
		Run()
	if err != nil {
		log.Fatal(err)
	}
}
