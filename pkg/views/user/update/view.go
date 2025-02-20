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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func UpdateUserProfileView(updateView *models.UserProfile) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Email").
				Value(&updateView.Email).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("email cannot be empty or only spaces")
					}
					if isVaild := utils.ValidateEmail(str); !isVaild {
						return errors.New("please enter correct email format")
					}
					return nil
				}),
			huh.NewInput().
				Title("Realname").
				Value(&updateView.Realname).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("real name cannot be empty")
					}
					if isValid := utils.ValidateName(str); !isValid {
						return errors.New("realname with illegal length or contains illegal characters")
					}
					return nil
				}),
			huh.NewInput().
				Title("Comment").
				Value(&updateView.Comment),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
