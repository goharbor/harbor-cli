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
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func UpdateScannerView(scanner *models.ScannerRegistration) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&scanner.Name).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("Description").
				Value(&scanner.Description),
			huh.NewSelect[string]().
				Title("Authentication Approach").
				Options(
					huh.NewOption("None", "None"),
					huh.NewOption("Basic", "Basic"),
					huh.NewOption("Bearer", "Bearer"),
					huh.NewOption("API-Key", "X-ScannerAdapter-API-Key"),
				).
				Value(&scanner.Auth),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}

	switch scanner.Auth {
	case "Basic":
		var username, password string
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Username").
					Value(&username).
					Validate(func(str string) error {
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
						if err := utils.ValidatePassword(str); err != nil {
							return err
						}
						return nil
					}),
			),
		).WithTheme(theme).Run()
		if err != nil {
			log.Fatal(err)
		}
		scanner.AccessCredential = username + ":" + password
	case "Bearer", "X-ScannerAdapter-API-Key":
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Token / API Key").
					Value(&scanner.AccessCredential).
					Validate(huh.ValidateNotEmpty()),
			),
		).WithTheme(theme).Run()
		if err != nil {
			log.Fatal(err)
		}
		if scanner.Auth == "Bearer" {
			scanner.AccessCredential = "Bearer: " + scanner.AccessCredential
		} else {
			scanner.AccessCredential = "APIKey: " + scanner.AccessCredential
		}
	}

	var url string = scanner.URL.String()
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Scanner Adapter URL").
				Value(&url).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("url cannot be empty")
					}
					formattedUrl := utils.FormatUrl(str)
					if err := utils.ValidateURL(formattedUrl); err != nil {
						return err
					}
					return nil
				}),
			huh.NewSelect[bool]().
				Title("Disable ?").
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(scanner.Disabled),
			huh.NewSelect[bool]().
				Title("Skip Certificate Verification ?").
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(scanner.SkipCertVerify),
			huh.NewSelect[bool]().
				Title("Use Internal Registry Address ?").
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(scanner.UseInternalAddr),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
	scanner.URL = strfmt.URI(utils.FormatUrl(url))
}
