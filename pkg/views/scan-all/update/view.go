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
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func UpdateSchedule(cron *string) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the cron expression").
				Description("Standard 6-field cron format: second minute hour day-of-month month day-of-week").
				Placeholder("0 0 0 * * *"). // Daily at midnight with seconds
				Value(cron).
				Validate(utils.ValidateCronExpression),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
