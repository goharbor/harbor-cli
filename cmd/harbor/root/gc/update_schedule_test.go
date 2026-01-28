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

package gc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCron_Empty(t *testing.T) {
	result, err := validateCron("")
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "cron expression cannot be empty")
}

func TestValidateCron_Valid6Field(t *testing.T) {
	cron := "0 0 12 * * *"
	result, err := validateCron(cron)
	assert.NoError(t, err)
	assert.Equal(t, cron, result)
}

func TestValidateCron_Valid6FieldWithSeconds(t *testing.T) {
	cron := "30 15 10 * * 1"
	result, err := validateCron(cron)
	assert.NoError(t, err)
	assert.Equal(t, cron, result)
}

func TestValidateCron_5FieldConversion(t *testing.T) {
	fiveFieldCron := "0 12 * * 1"
	expectedResult := "0 " + fiveFieldCron

	result, err := validateCron(fiveFieldCron)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestValidateCron_TooFewFields(t *testing.T) {
	cron := "0 12 *"
	result, err := validateCron(cron)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "6-field cron format")
}

func TestValidateCron_TooManyFields(t *testing.T) {
	cron := "0 0 12 * * * *"
	result, err := validateCron(cron)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "too many fields")
}

func TestValidateCron_SingleField(t *testing.T) {
	cron := "0"
	result, err := validateCron(cron)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "6-field cron format")
}

func TestValidateCron_Weekly(t *testing.T) {
	cron := "0 0 * * 0"
	result, err := validateCron(cron)
	assert.NoError(t, err)
	assert.Equal(t, "0 0 0 * * 0", result)
}

func TestValidateCron_Monthly(t *testing.T) {
	cron := "0 0 1 * *"
	result, err := validateCron(cron)
	assert.NoError(t, err)
	assert.Equal(t, "0 0 0 1 * *", result)
}

func TestValidateCron_Daily(t *testing.T) {
	cron := "0 0 * * *"
	result, err := validateCron(cron)
	assert.NoError(t, err)
	assert.Equal(t, "0 0 0 * * *", result)
}

func TestValidateCron_Hourly(t *testing.T) {
	cron := "* * * * *"
	result, err := validateCron(cron)
	assert.NoError(t, err)
	assert.Equal(t, "0 * * * * *", result)
}

func TestValidateCron_NonStandardSeconds(t *testing.T) {
	cron := "30 * * * * *"
	result, err := validateCron(cron)
	assert.NoError(t, err)
	assert.Equal(t, cron, result)
}

func TestValidateCron_Midnight(t *testing.T) {
	cron := "0 0 0 * * *"
	result, err := validateCron(cron)
	assert.NoError(t, err)
	assert.Equal(t, cron, result)
}

func TestValidScheduleTypes(t *testing.T) {
	assert.True(t, validScheduleTypes["None"])
	assert.True(t, validScheduleTypes["Hourly"])
	assert.True(t, validScheduleTypes["Daily"])
	assert.True(t, validScheduleTypes["Weekly"])
	assert.True(t, validScheduleTypes["Custom"])
	assert.False(t, validScheduleTypes["Monthly"])
	assert.False(t, validScheduleTypes["Yearly"])
	assert.False(t, validScheduleTypes[""])
}

func TestValidateCron_FieldsCount(t *testing.T) {
	testCases := []struct {
		name        string
		cron        string
		expectError bool
	}{
		{"0 fields", "", true},
		{"1 field", "0", true},
		{"2 fields", "0 0", true},
		{"3 fields", "0 0 0", true},
		{"4 fields", "0 0 0 0", true},
		{"5 fields - daily", "0 0 * * *", false},
		{"5 fields - weekly", "0 0 * * 0", false},
		{"5 fields - monthly", "0 0 1 * *", false},
		{"6 fields", "0 0 0 * * *", false},
		{"7 fields", "0 0 0 * * * 0", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validateCron(tc.cron)
			if tc.expectError {
				assert.Error(t, err, "Expected error for cron: %s", tc.cron)
				assert.Equal(t, "", result)
			} else {
				assert.NoError(t, err, "Expected no error for cron: %s", tc.cron)
				if result != "" {
					fields := strings.Fields(result)
					assert.Len(t, fields, 6, "Result should have 6 fields")
				}
			}
		})
	}
}
