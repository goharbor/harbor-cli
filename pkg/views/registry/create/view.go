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
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

// struct to hold registry options
type RegistryOption struct {
	ID   string
	Name string
}

func CreateRegistryView(createView *api.CreateRegView) {
	registries, _ := api.GetRegistryProviders()

	// Initialize a slice to hold registry options
	var registryOptions []RegistryOption

	// Iterate over registries to populate registryOptions
	for i, registry := range registries {
		registryOptions = append(registryOptions, RegistryOption{
			ID:   strconv.FormatInt(int64(i), 10),
			Name: registry,
		})
	}

	// Initialize a slice to hold select options
	var registrySelectOptions []huh.Option[string]

	// Iterate over registryOptions to populate registrySelectOptions
	for _, option := range registryOptions {
		registrySelectOptions = append(
			registrySelectOptions,
			huh.NewOption(option.Name, option.Name),
		)
	}

	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a Registry Provider").
				Value(&createView.Type).
				Options(registrySelectOptions...).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("registry provider cannot be empty or only spaces")
					}
					return nil
				}),

			huh.NewInput().
				Title("Name").
				Value(&createView.Name).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("name cannot be empty or only spaces")
					}
					if isValid := utils.ValidateRegistryName(str); !isValid {
						return errors.New("please enter the correct name format")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description),
			huh.NewInput().
				Title("URL").
				Value(&createView.URL).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("url cannot be empty or only spaces")
					}
					formattedUrl := utils.FormatUrl(str)
					if err := utils.ValidateURL(formattedUrl); err != nil {
						return err
					}
					// Update the bound value to the normalized URL after successful validation.
					createView.URL = formattedUrl
					return nil
				}),
			huh.NewInput().
				Title("Access Key").
				Value(&createView.Credential.AccessKey),
			huh.NewInput().
				Title("Access Secret").
				Value(&createView.Credential.AccessSecret),
			huh.NewConfirm().
				Title("Verify Cert").
				Value(&createView.Insecure).
				Affirmative("yes").
				Negative("no"),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}
