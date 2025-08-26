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
package change

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type PasswordChangeView struct {
	OldPassword     string
	NewPassword     string
	ConfirmPassword string
}

func ChangePasswordView(view *PasswordChangeView) {
	theme := huh.ThemeCharm()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Old Password").
				EchoMode(huh.EchoModePassword).
				Value(&view.OldPassword).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("old password cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("New Password").
				EchoMode(huh.EchoModePassword).
				Value(&view.NewPassword).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("new password cannot be empty")
					}
					if err := utils.ValidatePassword(str); err != nil {
						return err
					}
					return nil
				}),
			huh.NewInput().
				Title("Confirm New Password").
				EchoMode(huh.EchoModePassword).
				Value(&view.ConfirmPassword).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("confirmation password cannot be empty")
					}
					if str != view.NewPassword {
						return errors.New("passwords do not match")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
