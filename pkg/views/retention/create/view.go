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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	ScopeSelectors   RetentionSelector
	TagSelectors     RetentionSelector
	KeepLatestPushed int64
	Cron             string
}

type RetentionSelector struct {
	Decoration string
	Pattern    string
}

// validatePattern checks if a pattern string is valid for retention policies.
// Valid patterns can contain:
// - Alphanumeric characters, hyphens, and underscores
// - Wildcards (*)
// - Commas to separate multiple patterns (no consecutive, leading, or trailing commas)
func validatePattern(pattern string) error {
	if strings.TrimSpace(pattern) == "" {
		return errors.New("pattern cannot be empty")
	}

	// Check for leading or trailing commas
	if strings.HasPrefix(pattern, ",") || strings.HasSuffix(pattern, ",") {
		return errors.New("pattern cannot start or end with a comma")
	}

	// Check for consecutive commas
	if strings.Contains(pattern, ",,") {
		return errors.New("pattern cannot contain consecutive commas")
	}

	// Split by comma and validate each segment
	segments := strings.Split(pattern, ",")
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9\*\-_]+$`)

	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			return errors.New("pattern segments cannot be empty")
		}
		if !validPattern.MatchString(segment) {
			return errors.New("pattern contains invalid characters; only alphanumeric, *, -, and _ are allowed")
		}
	}

	return nil
}

// validateCron checks if a cron expression string is valid.
// Requires a standard 6-field cron format: second minute hour day-of-month month day-of-week
func validateCron(cronExpr string) error {
	if strings.TrimSpace(cronExpr) == "" {
		return errors.New("cron expression cannot be empty")
	}

	fields := strings.Fields(cronExpr)
	if len(fields) != 6 {
		if len(fields) == 5 {
			return fmt.Errorf("you entered a 5-field cron expression, but Harbor requires 6 fields (with seconds)\n"+
				"Please add a seconds field at the beginning. For example: '0 %s'", cronExpr)
		}
		return fmt.Errorf("harbor requires exactly 6 fields in cron expressions (seconds minute hour day month weekday), got %d", len(fields))
	}

	// Regex pattern for validating each field's range
	// Seconds: 0-59, Minutes: 0-59, Hours: 0-23, Day: 1-31, Month: 1-12, Weekday: 0-6
	cronRegex := regexp.MustCompile(`^(\*|[0-9]|[1-5][0-9]|\*/[0-9]+) (\*|[0-9]|[1-5][0-9]|\*/[0-9]+) (\*|[0-9]|1[0-9]|2[0-3]|\*/[0-9]+) (\*|[1-9]|[12][0-9]|3[01]|\*/[0-9]+) (\*|[1-9]|1[0-2]|\*/[0-9]+) (\*|[0-6]|\*/[0-9]+)$`)
	if !cronRegex.MatchString(cronExpr) {
		return errors.New("invalid cron expression format\n" +
			"Required format: second minute hour day-of-month month day-of-week\n" +
			"Examples:\n" +
			"  0 0 0 * * *    - Daily at midnight\n" +
			"  0 0 */6 * * *  - Every 6 hours\n" +
			"  0 0 0 * * 0    - Weekly on Sunday at midnight")
	}
	return nil
}

func CreateRetentionView(createView *CreateView) {
	keepLatestRaw := "10"
	if createView.KeepLatestPushed == 0 {
		createView.KeepLatestPushed = 10
	} else {
		keepLatestRaw = strconv.FormatInt(createView.KeepLatestPushed, 10)
	}
	if createView.Cron == "" {
		createView.Cron = "0 0 0 * * *"
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("For the repositories").
				Options(
					huh.NewOption("matching", "repoMatches"),
					huh.NewOption("excluding", "repoExcludes"),
				).
				Value(&createView.ScopeSelectors.Decoration).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("repository decoration cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Repository selector pattern").
				Description("Enter **, repo*, or comma-separated patterns").
				Value(&createView.ScopeSelectors.Pattern).
				Validate(func(str string) error {
					return validatePattern(str)
				}),
			huh.NewSelect[string]().
				Title("For the tags").
				Options(
					huh.NewOption("matching", "matches"),
					huh.NewOption("excluding", "excludes"),
				).
				Value(&createView.TagSelectors.Decoration).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("tag decoration cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Tag selector pattern").
				Description("Enter **, v*, or comma-separated patterns").
				Value(&createView.TagSelectors.Pattern).
				Validate(func(str string) error {
					return validatePattern(str)
				}),
			huh.NewInput().
				Title("Keep latest pushed artifacts").
				Description("Default 10").
				Value(&keepLatestRaw).
				Validate(func(v string) error {
					n, err := strconv.ParseInt(v, 10, 64)
					if err != nil || n <= 0 {
						return errors.New("keep latest value must be greater than 0")
					}
					return nil
				}),
			huh.NewInput().
				Title("Cron schedule").
				Description("Default 0 0 0 * * *").
				Value(&createView.Cron).
				Validate(func(str string) error {
					return validateCron(str)
				}),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		log.Fatal(err)
	}

	parsed, err := strconv.ParseInt(keepLatestRaw, 10, 64)
	if err != nil || parsed <= 0 {
		createView.KeepLatestPushed = 10
		return
	}

	createView.KeepLatestPushed = parsed
}
