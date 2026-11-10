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
package scan_all

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeCron(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "5-field cron gets seconds prepended",
			input:    "30 2 * * *",
			expected: "0 30 2 * * *",
		},
		{
			name:     "6-field cron is unchanged",
			input:    "0 30 2 * * *",
			expected: "0 30 2 * * *",
		},
		{
			name:     "extra whitespace is collapsed",
			input:    "  30  2 * * * ",
			expected: "0 30 2 * * *",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalizeCron(tt.input))
		})
	}
}

func TestValidateCron(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "6-field cron is valid", input: "0 30 2 * * *", wantErr: false},
		{name: "4-field cron is invalid", input: "2 * * *", wantErr: true},
		{name: "7-field cron is invalid", input: "0 30 2 * * * *", wantErr: true},
		{name: "empty cron is invalid", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCron(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCron_FiveFieldNormalized(t *testing.T) {
	// The documented behavior: a 5-field cron expression is converted to
	// Harbor's 6-field format by prepending the seconds field.
	err := validateCron(normalizeCron("30 2 * * *"))
	assert.NoError(t, err)
}
