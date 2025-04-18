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

	"github.com/charmbracelet/huh"
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
				Validate(func(str string) error {
					if str == "" {
						return errors.New("name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("description cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Authentication Approach").
				Options(
					huh.NewOption("None", "None"),
					huh.NewOption("Basic", "Basic"),
					huh.NewOption("Bearer", "Bearer"),
					huh.NewOption("X-ScannerAdapter-API-Key", "X-ScannerAdapter-API-Key"),
				).
				Value(&createView.Auth),
			huh.NewInput().
				Title("Access Credential").
				Value(&createView.AccessCredential).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("access credential cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("URL").
				Value(&createView.URL).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("url cannot be empty")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Disable?").
				Affirmative("Yes").
				Negative("No").
				Value(&createView.Disabled),
			huh.NewConfirm().
				Title("Skip Certificate Verification?").
				Affirmative("Yes").
				Negative("No").
				Value(&createView.SkipCertVerify),
			huh.NewConfirm().
				Title("Use Internal Registry Address?").
				Affirmative("Yes").
				Negative("No").
				Value(&createView.UseInternalAddr),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
