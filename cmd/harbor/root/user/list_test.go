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
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
					t.Errorf(`Expected logs to contain "No users found" but got: %s`, logs)
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

type MockUserLister struct {
	numberOfUsersforTesting int64
	usersForTesting         []*models.UserResp
	expectAuthError         bool
}

func (m *MockUserLister) populateUsers() []*models.UserResp {
	users := make([]*models.UserResp, 0, m.numberOfUsersforTesting)
	for i := 0; i < int(m.numberOfUsersforTesting); i++ {
		user := &models.UserResp{
			UserID: int64(i + 1),
		}
		users = append(users, user)
	}
	m.usersForTesting = users
	return users
}
func (m *MockUserLister) UserList(opts ...api.ListFlags) (*user.ListUsersOK, error) {
	if m.expectAuthError {
		return nil, fmt.Errorf("403")
	}
	res := &user.ListUsersOK{}
	if len(opts) == 0 {
		return res, nil
	}
	listFlags := opts[0]
	page, pageSize := listFlags.Page, listFlags.PageSize
	users := m.populateUsers()
	lo, hi := max(pageSize*(page-1), 0), min(pageSize*page, m.numberOfUsersforTesting)
	if lo >= m.numberOfUsersforTesting {
		return res, nil
	}
	res.Payload = users[lo:hi]
	return res, nil
}
func TestGetUsers(t *testing.T) {
	usersAreEqual := func(u1, u2 []*models.UserResp) bool {
		if len(u1) != len(u2) {
			return false
		}
		slices.SortFunc(u1, func(a, b *models.UserResp) int {
			return cmp.Compare(a.UserID, b.UserID)
		})
		slices.SortFunc(u2, func(a, b *models.UserResp) int {
			return cmp.Compare(a.UserID, b.UserID)
		})
		for i := 0; i < len(u1); i++ {
			if u1[i].UserID != u2[i].UserID {
				return false
			}
		}
		return true
	}
	tests := []struct {
		name        string
		setup       func() (api.ListFlags, *MockUserLister)
		wantError   bool
		errContains string
	}{
		{
			name: "fetch specific page with valid page size",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     2,
					PageSize: 50,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 102,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError: false,
		},
		{
			name: "fetch all users with page size 0 (multiple pages)",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     1,
					PageSize: 0,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 250,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError: false,
		},
		{
			name: "fetch all users when total is exactly divisible by 100",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     1,
					PageSize: 0,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 200,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError: false,
		},
		{
			name: "fetch first page with page size 10",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     1,
					PageSize: 10,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 50,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError: false,
		},
		{
			name: "fetch last page with partial results",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     3,
					PageSize: 10,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 25,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError: false,
		},
		{
			name: "fetch page beyond available data returns empty",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     10,
					PageSize: 10,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 5,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError: false,
		},
		{
			name: "fetch with maximum allowed page size 100",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     1,
					PageSize: 100,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 150,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError: false,
		},
		{
			name: "fetch with zero users in database",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     1,
					PageSize: 10,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 0,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError: false,
		},
		{
			name: "page size exceeds maximum (101)",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     1,
					PageSize: 101,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 50,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError:   true,
			errContains: "page size should be greater than or equal to 0 and less than or equal to 100",
		},
		{
			name: "page size is negative",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     1,
					PageSize: -1,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 50,
					expectAuthError:         false,
				}
				return opts, m
			},
			wantError:   true,
			errContains: "page size should be greater than or equal to 0 and less than or equal to 100",
		},
		{
			name: "authentication error returns permission denied",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     1,
					PageSize: 10,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 50,
					expectAuthError:         true,
				}
				return opts, m
			},
			wantError:   true,
			errContains: "Permission denied",
		},
		{
			name: "authentication error during fetch all",
			setup: func() (api.ListFlags, *MockUserLister) {
				opts := api.ListFlags{
					Page:     1,
					PageSize: 0,
				}
				m := &MockUserLister{
					numberOfUsersforTesting: 50,
					expectAuthError:         true,
				}
				return opts, m
			},
			wantError:   true,
			errContains: "Permission denied",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, m := tt.setup()
			allUsers, err := GetUsers(opts, m)

			// Check if we expected an error but did not get one (or vice-versa)
			if (err != nil) != tt.wantError {
				t.Fatalf("GetUsers() error presence mismatch: got error %v, wantError %v", err, tt.wantError)
			}
			if tt.wantError && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Expected error to contain '%s', got '%s'", tt.errContains, err.Error())
			}
			if !tt.wantError {
				// check if the returned allUsers are correct according to our mock database of users (which is m.usersForTesting)
				if opts.PageSize == 0 {
					if !usersAreEqual(allUsers, m.usersForTesting) {
						t.Errorf("Expected all of the users to be returned")
					}
				} else {
					requiredPage, requiredPageSize := opts.Page, opts.PageSize
					start := max(requiredPageSize*(requiredPage-1), 0)
					end := min(requiredPageSize*requiredPage, m.numberOfUsersforTesting)

					if start >= m.numberOfUsersforTesting {
						if len(allUsers) != 0 {
							t.Errorf("Expected empty result for page beyond data, got %d users", len(allUsers))
						}
					} else {
						if !usersAreEqual(allUsers, m.usersForTesting[start:end]) {
							t.Errorf("Expected different set of users")
						}
					}
				}
			}
		})
	}
}
func TestUserListCmd(t *testing.T) {
	cmd := UserListCmd()

	assert.Equal(t, "list", cmd.Use)
	assert.Equal(t, "List users", cmd.Short)
	assert.Contains(t, cmd.Aliases, "ls")

	pageFlag := cmd.Flags().Lookup("page")
	assert.NotNil(t, pageFlag)
	assert.Equal(t, "1", pageFlag.DefValue)

	pageSizeFlag := cmd.Flags().Lookup("page-size")
	assert.NotNil(t, pageSizeFlag)
	assert.Equal(t, "0", pageSizeFlag.DefValue)

	queryFlag := cmd.Flags().Lookup("query")
	assert.NotNil(t, queryFlag)

	sortFlag := cmd.Flags().Lookup("sort")
	assert.NotNil(t, sortFlag)

	assert.Equal(t, "p", pageFlag.Shorthand)
	assert.Equal(t, "n", pageSizeFlag.Shorthand)
	assert.Equal(t, "q", queryFlag.Shorthand)
	assert.Equal(t, "s", sortFlag.Shorthand)
}
