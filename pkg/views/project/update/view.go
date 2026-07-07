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
	"github.com/goharbor/harbor-cli/pkg/utils"
)

type UpdateView struct {
	Public       string
	StorageLimit string
	RegistryID   string
	AutoScan     *string
	PreventVul   *string
	Severity     *string
	ReuseSysCVE  *string
}

func validateValue(value *string) *string {
	defaultVal := "false"
	if value == nil {
		return &defaultVal
	}
	return value
}

func UpdateProjectView(view *UpdateView) error {
	theme := huh.ThemeCharm()
	groups := []*huh.Group{
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Make Project Public").
				Options(
					huh.NewOption("No", "false"),
					huh.NewOption("Yes", "true"),
				).
				Value(&view.Public),
			huh.NewInput().
				Title("Storage Limit").
				Value(&view.StorageLimit).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("storage limit cannot be empty")
					}
					return utils.ValidateStorageLimit(str)
				}),
			huh.NewSelect[string]().
				Title("Automatically scan images on push").
				Options(
					huh.NewOption("No", "false"),
					huh.NewOption("Yes", "true"),
				).
				Value(validateValue(view.AutoScan)),
			huh.NewSelect[string]().
				Title("Prevent vulnerable images from running").
				Options(
					huh.NewOption("No", "false"),
					huh.NewOption("Yes", "true"),
				).
				Value(validateValue(view.PreventVul)),
			huh.NewSelect[string]().
				Title("Vulnerability severity threshold").
				Options(
					huh.NewOption("None", "none"),
					huh.NewOption("Low", "low"),
					huh.NewOption("Medium", "medium"),
					huh.NewOption("High", "high"),
					huh.NewOption("Critical", "critical"),
				).
				Value(validateValue(view.Severity)),
			huh.NewSelect[string]().
				Title("Reuse system CVE allowlist").
				Options(
					huh.NewOption("No", "false"),
					huh.NewOption("Yes", "true"),
				).
				Value(validateValue(view.ReuseSysCVE)),
		),
	}

	if view.RegistryID != "" {
		groups = append(groups, huh.NewGroup(
			huh.NewInput().
				Title("Registry ID").
				Value(&view.RegistryID),
		))
	}

	return huh.NewForm(groups...).WithTheme(theme).Run()
}
