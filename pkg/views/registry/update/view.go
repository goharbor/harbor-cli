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

func UpdateRegistryView(updateView *models.Registry) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Provider").
				Value(&updateView.Type).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("provider cannot be empty or only spaces")
					}
					return nil
				}),
			huh.NewInput().
				Title("Name").
				Value(&updateView.Name).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("name cannot be empty or only spaces")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&updateView.Description),
			huh.NewInput().
				Title("URL").
				Value(&updateView.URL).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("url cannot be empty or only spaces")
					}
					formattedUrl := utils.FormatUrl(str)
					if err := utils.ValidateURL(formattedUrl); err != nil {
						return err
					}
					updateView.URL = formattedUrl
					return nil
				}),
			huh.NewInput().
				Title("Access ID").
				Value(&updateView.Credential.AccessKey),
			huh.NewInput().
				Title("Access Secret").
				EchoMode(huh.EchoModePassword).
				Description("Replace the Access Secret to the real one").
				Value(&updateView.Credential.AccessSecret),
			huh.NewConfirm().
				Title("Verify Cert").
				Value(&updateView.Insecure).
				Affirmative("yes").
				Negative("no"),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
