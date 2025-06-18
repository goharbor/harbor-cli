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
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
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
				Validate(validateCronExpression),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}

func validateCronExpression(cron string) error {
	if cron == "" {
		return errors.New("cron expression cannot be empty")
	}
	fields := strings.Fields(cron)
	if len(fields) != 6 {
		if len(fields) == 5 {
			return fmt.Errorf("you entered a 5-field cron expression, but Harbor requires 6 fields (with seconds)\n"+
				"Please add a seconds field at the beginning. For example: '0 %s'", cron)
		}
		return fmt.Errorf("harbor requires exactly 6 fields in cron expressions (seconds minute hour day month weekday), got %d", len(fields))
	}
	cronRegex := regexp.MustCompile(`^(\*|[0-9]|[1-5][0-9]|\*/[0-9]+) (\*|[0-9]|[1-5][0-9]|\*/[0-9]+) (\*|[0-9]|1[0-9]|2[0-3]|\*/[0-9]+) (\*|[1-9]|[12][0-9]|3[01]|\*/[0-9]+) (\*|[1-9]|1[0-2]|\*/[0-9]+) (\*|[0-6]|\*/[0-9]+)$`)
	if !cronRegex.MatchString(cron) {
		return errors.New("invalid cron expression format\n" +
			"Examples:\n" +
			"  0 0 0 * * *    - Daily at midnight\n" +
			"  0 0 */6 * * *  - Every 6 hours\n" +
			"  0 0 0 * * 0    - Weekly on Sunday at midnight")
	}
	return nil
}
