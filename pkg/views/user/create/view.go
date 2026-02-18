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
	Username        string
	Email           string
	Realname        string
	Comment         string
	Password        string
	ConfirmPassword string
}

func CreateUserView(createView *CreateView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("User Name").
				Value(&createView.Username).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("user name cannot be empty")
					}
					if isValid := utils.ValidateUserName(str); !isValid {
						return errors.New("username cannot contain special characters")
					}
					return nil
				}),
			huh.NewInput().
				Title("Email").
				Value(&createView.Email).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("email cannot be empty or only spaces")
					}
					if isValid := utils.ValidateEmail(str); !isValid {
						return errors.New("please enter correct email format")
					}
					return nil
				}),

			huh.NewInput().
				Title("First and Last Name").
				Value(&createView.Realname).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("real name cannot be empty")
					}
					if isValid := utils.ValidateFL(str); !isValid {
						return errors.New("please enter correct first and last name format, like `Bob Dylan`")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&createView.Password).
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
				Title("Confirm Password").
				EchoMode(huh.EchoModePassword).
				Value(&createView.ConfirmPassword).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("confirm password cannot be empty")
					}
					if str != createView.Password {
						return errors.New("passwords do not match")
					}
					return nil
				}),
			huh.NewInput().
				Title("Comment (optional)").
				Value(&createView.Comment),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
