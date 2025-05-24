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
	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	log "github.com/sirupsen/logrus"
)

func validateValue(value *string) *string {
	defaultVal := "false"
	if value == nil {
		return &defaultVal
	}
	return value
}

func UpdateProjectMetadataView(config *models.ProjectMetadata) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Make Project Public").
				Options(
					huh.NewOption("No", "false"),
					huh.NewOption("Yes", "true"),
				).
				Value(&config.Public),

			huh.NewSelect[string]().
				Title("Automatically scan images on push").
				Options(
					huh.NewOption("No", "false"),
					huh.NewOption("Yes", "true"),
				).
				Value(validateValue(config.AutoScan)),
			huh.NewSelect[string]().
				Title("Prevent vulnerable images from running").
				Options(
					huh.NewOption("No", "false"),
					huh.NewOption("Yes", "true"),
				).
				Value(validateValue(config.PreventVul)),
			huh.NewSelect[string]().
				Title("Vulnerability severity threshold").
				Options(
					huh.NewOption("None", "none"),
					huh.NewOption("Low", "low"),
					huh.NewOption("Medium", "medium"),
					huh.NewOption("High", "high"),
					huh.NewOption("Critical", "critical"),
				).
				Value(validateValue(config.Severity)),
			huh.NewSelect[string]().
				Title("Reuse system CVE allowlist").
				Options(
					huh.NewOption("No", "false"),
					huh.NewOption("Yes", "true"),
				).
				Value(validateValue(config.ReuseSysCVEAllowlist)),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
