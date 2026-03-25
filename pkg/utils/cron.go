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
	"strings"

	"github.com/sirupsen/logrus"
)

// ValidateCron validates and normalizes a cron expression for Harbor.
// Harbor requires 6-field cron format (seconds minute hour day month weekday).
// If a 5-field expression is provided, it prepends '0' for seconds.
func ValidateCron(cron string) (string, error) {
	if cron == "" {
		return "", errors.New("cron expression cannot be empty")
	}
	fields := strings.Fields(cron)
	if len(fields) < 6 {
		if len(fields) == 5 {
			logrus.Infof("Converting 5-field cron to 6-field by adding '0' for seconds")
			return fmt.Sprintf("0 %s", cron), nil
		}
		return "", fmt.Errorf("harbor requires 6-field cron format (seconds minute hour day month weekday)")
	}
	if len(fields) > 6 {
		return "", fmt.Errorf("too many fields in cron expression, expected 6 but got %d", len(fields))
	}
	return cron, nil
}
