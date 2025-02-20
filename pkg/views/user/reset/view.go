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
package reset

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type ResetView struct {
	OldPassword string
	NewPassword string
	Comfirm     string
}

func ResetUserView(resetView *ResetView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Old Password").
				EchoMode(huh.EchoModePassword).
				Value(&resetView.OldPassword).
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
				Title("New Password").
				EchoMode(huh.EchoModePassword).
				Value(&resetView.NewPassword).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("password cannot be empty or only spaces")
					}
					if resetView.OldPassword == resetView.NewPassword {
						return errors.New("new password is the same as the old one")
					}
					if err := utils.ValidatePassword(str); err != nil {
						return err
					}
					return nil
				}),
			huh.NewInput().
				Title("Comfirm Password").
				EchoMode(huh.EchoModePassword).
				Value(&resetView.Comfirm).
				Validate(func(str string) error {
					if resetView.Comfirm != resetView.NewPassword {
						return errors.New("passwords do not match, please try again")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
