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
package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// cronRegex validates a 6-field Harbor cron expression:
//
//	second minute hour day-of-month month day-of-week
//
// Each field accepts * (any), a numeric value in the allowed range, or */n
// (step notation). This mirrors the regex previously inlined in
// pkg/views/scan-all/update/view.go.
var cronRegex = regexp.MustCompile(
	`^(\*|[0-9]|[1-5][0-9]|\*/[0-9]+) ` +
		`(\*|[0-9]|[1-5][0-9]|\*/[0-9]+) ` +
		`(\*|[0-9]|1[0-9]|2[0-3]|\*/[0-9]+) ` +
		`(\*|[1-9]|[12][0-9]|3[01]|\*/[0-9]+) ` +
		`(\*|[1-9]|1[0-2]|\*/[0-9]+) ` +
		`(\*|[0-6]|\*/[0-9]+)$`,
)

// ValidateCronExpression checks that expr is a valid 6-field cron string
// as accepted by Harbor (seconds, minutes, hours, day-of-month, month, day-of-week).
// Returns nil on success, a descriptive error on failure.
//
// This function can be used directly as a huh form Validate callback because
// its signature matches func(string) error.
func ValidateCronExpression(expr string) error {
	if expr == "" {
		return errors.New("cron expression cannot be empty")
	}

	fields := strings.Fields(expr)
	switch {
	case len(fields) == 5:
		return fmt.Errorf(
			"you entered a 5-field cron expression, but Harbor requires 6 fields (with seconds)\n"+
				"Please add a seconds field at the beginning. For example: '0 %s'", expr,
		)
	case len(fields) < 6:
		return fmt.Errorf(
			"harbor requires exactly 6 fields in cron expressions (seconds minute hour day month weekday), got %d",
			len(fields),
		)
	case len(fields) > 6:
		return fmt.Errorf(
			"too many fields in cron expression, expected 6 but got %d",
			len(fields),
		)
	}

	if !cronRegex.MatchString(expr) {
		return errors.New("invalid cron expression format\n" +
			"Examples:\n" +
			"  0 0 0 * * *    - Daily at midnight\n" +
			"  0 0 */6 * * *  - Every 6 hours\n" +
			"  0 0 0 * * 0    - Weekly on Sunday at midnight")
	}

	return nil
}
