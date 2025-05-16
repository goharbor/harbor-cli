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
	Name             string
	Description      string
	Auth             string
	AccessCredential string
	URL              string
	Disabled         bool
	SkipCertVerify   bool
	UseInternalAddr  bool
}

func CreateScannerView(createView *CreateView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&createView.Name).
				Validate(huh.ValidateNotEmpty()),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description),
			huh.NewSelect[string]().
				Title("Authentication Approach").
				Options(
					huh.NewOption("None", "None"),
					huh.NewOption("Basic", "Basic"),
					huh.NewOption("Bearer", "Bearer"),
					huh.NewOption("API-Key", "X-ScannerAdapter-API-Key"),
				).
				Value(&createView.Auth),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}

	if createView.Auth == "Basic" {
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
		createView.AccessCredential = username + ":" + password
	} else if createView.Auth == "Bearer" || createView.Auth == "X-ScannerAdapter-API-Key" {
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Token / API Key").
					Value(&createView.AccessCredential).
					Validate(huh.ValidateNotEmpty()),
			),
		).WithTheme(theme).Run()
		if err != nil {
			log.Fatal(err)
		}
		if createView.Auth == "Bearer" {
			createView.AccessCredential = "Bearer: " + createView.AccessCredential
		} else {
			createView.AccessCredential = "APIKey: " + createView.AccessCredential
		}
	}
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("URL").
				Value(&createView.URL).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("server cannot be empty or only spaces")
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
				Value(&createView.Disabled),
			huh.NewSelect[bool]().
				Title("Skip Certificate Verification ?").
				Description("Skip TLS check (use for self-signed certs).").
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(&createView.SkipCertVerify),
			huh.NewSelect[bool]().
				Title("Use Internal Registry Address ?").
				Description("Use internal Harbor address (for in-cluster scanners).").
				Options(
					huh.NewOption("No", false),
					huh.NewOption("Yes", true),
				).
				Value(&createView.UseInternalAddr),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
	createView.URL = utils.FormatUrl(createView.URL)
}
