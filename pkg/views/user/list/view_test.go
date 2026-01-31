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
package list

import (
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func TestMakeUserRows(t *testing.T) {
	dateStr := "2023-01-01T12:00:00Z"
	testDate, err := strfmt.ParseDateTime(dateStr)
	if err != nil {
		t.Fatalf("failed to parse date %q: %v", dateStr, err)
	}
	expectedTimeStr, err := utils.FormatCreatedTime(dateStr)
	if err != nil {
		t.Fatalf("failed to format created time %q: %v", dateStr, err)
	}
	tests := []struct {
		name     string
		setup    func() []*models.UserResp
		expected [][]string
	}{
		{
			name: "Number of users non-zero",
			setup: func() []*models.UserResp {
				return []*models.UserResp{
					{
						UserID:       1,
						Username:     "testUser1",
						Email:        "test1@domain.com",
						SysadminFlag: true,
						Realname:     "Test1",
						CreationTime: testDate,
					},
					{
						UserID:       2,
						Username:     "testUser2",
						Email:        "test2@domain.com",
						SysadminFlag: false,
						Realname:     "Test2",
						CreationTime: testDate,
					},
					{
						UserID:       3,
						Username:     "testUser3",
						Email:        "test3@domain.com",
						SysadminFlag: false,
						Realname:     "Test3",
						CreationTime: testDate,
					},
				}
			},
			expected: [][]string{
				{"1", "testUser1", "Yes", "test1@domain.com", expectedTimeStr},
				{"2", "testUser2", "No", "test2@domain.com", expectedTimeStr},
				{"3", "testUser3", "No", "test3@domain.com", expectedTimeStr},
			},
		},
		{
			name: "No users",
			setup: func() []*models.UserResp {
				return []*models.UserResp{}
			},
			expected: [][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := tt.setup()
			rows := MakeUserRows(users)
			if len(tt.expected) != len(rows) {
				t.Fatalf("MakeUserRows returned %d rows for %d users", len(rows), len(users))
			}
			for i := 0; i < len(rows); i++ {
				if len(rows[i]) != len(tt.expected[i]) {
					t.Errorf("Row %d: expected %d columns, got %d", i, len(tt.expected[i]), len(rows[i]))
					continue
				}
				for j := 0; j < len(rows[i]); j++ {
					if rows[i][j] != tt.expected[i][j] {
						t.Errorf("Row %d, Column %d: expected '%s', but got '%s'",
							i, j, tt.expected[i][j], rows[i][j])
					}
				}
			}
		})
	}
}
