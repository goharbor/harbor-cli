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

package user

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.yaml.in/yaml/v4"
)

func TestPrintUsers(t *testing.T) {
	testDate, _ := strfmt.ParseDateTime("2023-01-01T12:00:00Z")
	testUsers := func() []*models.UserResp {
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
	}
	tests := []struct {
		name         string
		setup        func() []*models.UserResp
		outputFormat string
	}{
		{
			name: "Number of users not zero and output format is json",
			setup: func() []*models.UserResp {
				users := testUsers()
				return users
			},
			outputFormat: "json",
		},
		{
			name: "Number of users not zero and output format yaml",
			setup: func() []*models.UserResp {
				users := testUsers()
				return users
			},
			outputFormat: "yaml",
		},
		{
			name: "Number of users not zero and output format default",
			setup: func() []*models.UserResp {
				users := testUsers()
				return users
			},
			outputFormat: "",
		},
		{
			name: "Number of users is zero",
			setup: func() []*models.UserResp {
				users := []*models.UserResp{}
				return users
			},
			outputFormat: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allUsers := tt.setup()

			var logBuf, contentBuf bytes.Buffer
			originalLogOutput := log.StandardLogger().Out
			log.SetOutput(&logBuf)
			defer log.SetOutput(originalLogOutput)

			originalFormatFlag := viper.GetString("output-format")
			viper.Set("output-format", tt.outputFormat)
			defer viper.Set("output-format", originalFormatFlag)

			if err := PrintUsers(&contentBuf, allUsers); err != nil {
				t.Fatalf("PrintUsers() returned error: %v", err)
			}

			logs := logBuf.String()

			switch {
			case len(allUsers) == 0:
				if !strings.Contains(logs, "No users found") {
					t.Errorf(`Expected logs to contain "No user found" but got: %s`, logs)
				}
			case tt.outputFormat == "json":
				if contentBuf.Len() == 0 {
					t.Fatal("Expected JSON output, but buffer was empty")
				}
				var decodedUsers []*models.UserResp
				if err := json.Unmarshal(contentBuf.Bytes(), &decodedUsers); err != nil {
					t.Fatalf("Output is not valid JSON: %v. Output:\n%s", err, contentBuf.String())
				}
				if len(decodedUsers) != len(allUsers) {
					t.Errorf("Expected %d users in JSON, got %d", len(allUsers), len(decodedUsers))
				}
				if len(decodedUsers) > 0 {
					if decodedUsers[0].Username != allUsers[0].Username {
						t.Errorf("Expected username '%s', got '%s'", allUsers[0].Username, decodedUsers[0].Username)
					}
					if decodedUsers[0].SysadminFlag != allUsers[0].SysadminFlag {
						t.Errorf("Expected SysadminFlag to be %v, got %v", allUsers[0].SysadminFlag, decodedUsers[0].SysadminFlag)
					}
				}
			case tt.outputFormat == "yaml":
				if contentBuf.Len() == 0 {
					t.Fatal("Expected YAML output, but buffer was empty")
				}
				var decodedUsers []*models.UserResp
				if err := yaml.Unmarshal(contentBuf.Bytes(), &decodedUsers); err != nil {
					t.Fatalf("Output is not valid YAML: %v. Output:\n%s", err, contentBuf.String())
				}
				if len(decodedUsers) != len(allUsers) {
					t.Errorf("Expected %d users in YAML, got %d", len(allUsers), len(decodedUsers))
				}
				if len(decodedUsers) > 0 {
					if decodedUsers[0].Username != allUsers[0].Username {
						t.Errorf("Expected username '%s', got '%s'", allUsers[0].Username, decodedUsers[0].Username)
					}
					if decodedUsers[0].SysadminFlag != allUsers[0].SysadminFlag {
						t.Errorf("Expected SysadminFlag to be %v, got %v", allUsers[0].SysadminFlag, decodedUsers[0].SysadminFlag)
					}
				}
			default:
				if contentBuf.Len() == 0 {
					t.Fatal("Expected TUI table output, but buffer was empty. Did you pass 'w' to ListUsers?")
				}
				output := contentBuf.String()
				if !strings.Contains(output, "ID") || !strings.Contains(output, "Name") || !strings.Contains(output, "Administrator") {
					t.Error("Expected table output to contain headers 'ID', 'Name' and 'Administrator among other headers")
				}
				if !strings.Contains(output, "testUser1") {
					t.Errorf("Expected table to contain username 'testUser1'")
				}
			}
		})
	}
}
