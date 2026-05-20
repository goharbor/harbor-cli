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
package utils_test

import (
	"testing"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidateCronExpression(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// ---------- valid expressions ----------
		{
			name:    "daily at midnight",
			input:   "0 0 0 * * *",
			wantErr: false,
		},
		{
			name:    "every 6 hours",
			input:   "0 0 */6 * * *",
			wantErr: false,
		},
		{
			name:    "weekly on Sunday at midnight",
			input:   "0 0 0 * * 0",
			wantErr: false,
		},
		{
			name:    "all wildcards",
			input:   "* * * * * *",
			wantErr: false,
		},
		{
			name:    "specific time with step second",
			input:   "*/30 0 12 * * *",
			wantErr: false,
		},
		{
			name:    "first day of month at 3am",
			input:   "0 0 3 1 * *",
			wantErr: false,
		},
		{
			name:    "max values",
			input:   "59 59 23 31 12 6",
			wantErr: false,
		},

		// ---------- empty / blank ----------
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},

		// ---------- wrong field count ----------
		{
			name:    "5-field classic cron (no seconds)",
			input:   "0 0 * * *",
			wantErr: true,
		},
		{
			name:    "4 fields",
			input:   "0 0 * *",
			wantErr: true,
		},
		{
			name:    "7 fields (too many)",
			input:   "0 0 0 * * * extra",
			wantErr: true,
		},

		// ---------- bad field values ----------
		{
			name:    "hour out of range (25)",
			input:   "0 0 25 * * *",
			wantErr: true,
		},
		{
			name:    "month out of range (13)",
			input:   "0 0 0 * 13 *",
			wantErr: true,
		},
		{
			name:    "day-of-week out of range (7)",
			input:   "0 0 0 * * 7",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			input:   "abc def ghi * * *",
			wantErr: true,
		},
		{
			name:    "single field only",
			input:   "*",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := utils.ValidateCronExpression(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
