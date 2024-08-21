// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package login

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type LoginView struct {
	Server   string
	Username string
	Password string
	Name     string
}

func CreateView(loginView *LoginView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server").
				Value(&loginView.Server).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("server cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("User Name").
				Value(&loginView.Username).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("username cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&loginView.Password).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("password cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Name of Credential").
				Value(&loginView.Name).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("credential name cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}

}
