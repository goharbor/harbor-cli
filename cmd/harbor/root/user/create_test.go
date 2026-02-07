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
	"fmt"
	"reflect"
	"slices"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/goharbor/harbor-cli/pkg/views/user/create"
	"github.com/stretchr/testify/assert"
)

type MockUserCreator struct {
	usernames       map[string]struct{}
	emails          map[string]struct{}
	users           []*create.CreateView
	expectAuthError bool
}

func initMockUserCreator(expectAuthError bool) *MockUserCreator {
	return &MockUserCreator{
		usernames:       make(map[string]struct{}),
		emails:          make(map[string]struct{}),
		users:           []*create.CreateView{},
		expectAuthError: expectAuthError,
	}
}

/*
 FillUser simulates the interactive prompt that fills missing fields.
 Note: Username and Email must be provided in test cases since they're
 required fields that would be prompted interactively in real usage.
 Only Realname and Password are filled here to test the FillUser path.
*/

func (m *MockUserCreator) FillUser(opts *create.CreateView) {
	randomStr := "Random string 1234"
	if opts.Realname == "" {
		opts.Realname = randomStr
	}
	if opts.Password == "" {
		opts.Password = randomStr
	}
}

func (m *MockUserCreator) UserCreate(opts create.CreateView) error {
	if m.expectAuthError {
		return fmt.Errorf("403")
	}
	if opts.Email == "" || opts.Realname == "" || opts.Password == "" || opts.Username == "" {
		return fmt.Errorf("missing required fields")
	}
	_, foundUsername := m.usernames[opts.Username]
	_, foundEmail := m.emails[opts.Email]
	if foundUsername || foundEmail {
		return fmt.Errorf("user %s or email %s already exists", opts.Username, opts.Email)
	}
	m.usernames[opts.Username] = struct{}{}
	m.emails[opts.Email] = struct{}{}
	m.users = append(m.users, &create.CreateView{
		Username: opts.Username,
		Email:    opts.Email,
		Comment:  opts.Comment,
		Realname: opts.Realname,
		Password: opts.Password,
	})
	return nil
}
func TestCreateUsers(t *testing.T) {
	usersAreEqual := func(u1, u2 []*create.CreateView) bool {
		if len(u1) != len(u2) {
			return false
		}
		u1Copy := make([]*create.CreateView, len(u1))
		u2Copy := make([]*create.CreateView, len(u2))
		copy(u1Copy, u1)
		copy(u2Copy, u2)

		slices.SortFunc(u1Copy, func(a, b *create.CreateView) int {
			return cmp.Compare(a.Username, b.Username)
		})
		slices.SortFunc(u2Copy, func(a, b *create.CreateView) int {
			return cmp.Compare(a.Username, b.Username)
		})
		for i := 0; i < len(u1Copy); i++ {
			if !reflect.DeepEqual(u1Copy[i], u2Copy[i]) {
				return false
			}
		}
		return true
	}
	tests := []struct {
		name          string
		setup         func() ([]*create.CreateView, *MockUserCreator)
		expectedUsers []*create.CreateView
		expectedErr   string
	}{
		{
			name: "successfully create single user with all fields",
			setup: func() ([]*create.CreateView, *MockUserCreator) {
				views := []*create.CreateView{
					{
						Username: "testuser",
						Email:    "test@example.com",
						Realname: "Test User",
						Password: "TestPass123",
						Comment:  "Test comment",
					},
				}
				return views, initMockUserCreator(false)
			},
			expectedUsers: []*create.CreateView{
				{
					Username: "testuser",
					Email:    "test@example.com",
					Realname: "Test User",
					Password: "TestPass123",
					Comment:  "Test comment",
				},
			},
			expectedErr: "",
		},
		{
			name: "successfully create multiple users",
			setup: func() ([]*create.CreateView, *MockUserCreator) {
				views := []*create.CreateView{
					{Username: "user1", Email: "user1@example.com", Realname: "User One", Password: "Pass1", Comment: "First"},
					{Username: "user2", Email: "user2@example.com", Realname: "User Two", Password: "Pass2", Comment: "Second"},
					{Username: "user3", Email: "user3@example.com", Realname: "User Three", Password: "Pass3", Comment: "Third"},
				}
				return views, initMockUserCreator(false)
			},
			expectedUsers: []*create.CreateView{
				{Username: "user1", Email: "user1@example.com", Realname: "User One", Password: "Pass1", Comment: "First"},
				{Username: "user2", Email: "user2@example.com", Realname: "User Two", Password: "Pass2", Comment: "Second"},
				{Username: "user3", Email: "user3@example.com", Realname: "User Three", Password: "Pass3", Comment: "Third"},
			},
			expectedErr: "",
		},
		{
			name: "create user with missing fields triggers FillUser",
			setup: func() ([]*create.CreateView, *MockUserCreator) {
				views := []*create.CreateView{
					{
						Username: "testuser",
						Email:    "test@example.com",
					},
				}
				return views, initMockUserCreator(false)
			},
			expectedUsers: []*create.CreateView{
				{
					Username: "testuser",
					Email:    "test@example.com",
					Realname: "Random string 1234",
					Password: "Random string 1234",
				},
			},
			expectedErr: "",
		},
		{
			name: "permission denied error (403)",
			setup: func() ([]*create.CreateView, *MockUserCreator) {
				views := []*create.CreateView{
					{
						Username: "testuser",
						Email:    "test@example.com",
						Realname: "Test User",
						Password: "TestPass123",
					},
				}
				return views, initMockUserCreator(true)
			},
			expectedUsers: []*create.CreateView{},
			expectedErr:   "Permission denied",
		},
		{
			name: "duplicate username fails second user",
			setup: func() ([]*create.CreateView, *MockUserCreator) {
				views := []*create.CreateView{
					{Username: "sameuser", Email: "first@example.com", Realname: "First", Password: "Pass1"},
					{Username: "sameuser", Email: "second@example.com", Realname: "Second", Password: "Pass2"},
				}
				return views, initMockUserCreator(false)
			},
			expectedUsers: []*create.CreateView{
				{Username: "sameuser", Email: "first@example.com", Realname: "First", Password: "Pass1"},
			},
			expectedErr: "already exists",
		},
		{
			name: "duplicate email fails second user",
			setup: func() ([]*create.CreateView, *MockUserCreator) {
				views := []*create.CreateView{
					{Username: "user1", Email: "same@example.com", Realname: "First", Password: "Pass1"},
					{Username: "user2", Email: "same@example.com", Realname: "Second", Password: "Pass2"},
				}
				return views, initMockUserCreator(false)
			},
			expectedUsers: []*create.CreateView{
				{Username: "user1", Email: "same@example.com", Realname: "First", Password: "Pass1"},
			},
			expectedErr: "already exists",
		},
		{
			name: "create user with empty comment succeeds",
			setup: func() ([]*create.CreateView, *MockUserCreator) {
				views := []*create.CreateView{
					{
						Username: "testuser",
						Email:    "test@example.com",
						Realname: "Test User",
						Password: "TestPass123",
						Comment:  "",
					},
				}
				return views, initMockUserCreator(false)
			},
			expectedUsers: []*create.CreateView{
				{
					Username: "testuser",
					Email:    "test@example.com",
					Realname: "Test User",
					Password: "TestPass123",
					Comment:  "",
				},
			},
			expectedErr: "",
		},
		{
			name: "no users to create",
			setup: func() ([]*create.CreateView, *MockUserCreator) {
				views := []*create.CreateView{}
				return views, initMockUserCreator(false)
			},
			expectedUsers: []*create.CreateView{},
			expectedErr:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			originalLogOutput := log.StandardLogger().Out
			log.SetOutput(&buf)
			defer log.SetOutput(originalLogOutput)

			opts, m := tt.setup()

			for _, opt := range opts {
				CreateUser(opt, m)
			}
			logs := buf.String()
			if tt.expectedErr != "" {
				assert.Contains(t, logs, tt.expectedErr, "Expected error logs to contain %s but got %s", tt.expectedErr, logs)
			} else {
				assert.Empty(t, logs, "Expected no error logs but got: %s", logs)
			}
			if !usersAreEqual(m.users, tt.expectedUsers) {
				t.Errorf("Users mismatch.\nExpected: %+v\nGot: %+v", tt.expectedUsers, m.users)
			}
		})
	}
}

func TestIsUnauthorizedError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error returns false",
			err:      nil,
			expected: false,
		},
		{
			name:     "error containing 403 returns true",
			err:      fmt.Errorf("403 Forbidden"),
			expected: true,
		},
		{
			name:     "error without 403 returns false",
			err:      fmt.Errorf("404 Not Found"),
			expected: false,
		},
		{
			name:     "error with 500 returns false",
			err:      fmt.Errorf("500 Internal Server Error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isUnauthorizedError(tt.err)
			assert.Equal(t, tt.expected, result, "isUnauthorizedError(%v) = %v, want %v", tt.err, result, tt.expected)
		})
	}
}

func TestUserCreateCmd(t *testing.T) {
	cmd := UserCreateCmd()

	assert.Equal(t, "create", cmd.Use)
	assert.Equal(t, "create user", cmd.Short)
	assert.NotNil(t, cmd.Args, "Args validator should be set")

	emailFlag := cmd.Flags().Lookup("email")
	assert.NotNil(t, emailFlag)
	assert.Equal(t, "", emailFlag.DefValue)

	realnameFlag := cmd.Flags().Lookup("realname")
	assert.NotNil(t, realnameFlag)
	assert.Equal(t, "", realnameFlag.DefValue)

	commentFlag := cmd.Flags().Lookup("comment")
	assert.NotNil(t, commentFlag)
	assert.Equal(t, "", commentFlag.DefValue)

	passwordFlag := cmd.Flags().Lookup("password")
	assert.NotNil(t, passwordFlag)
	assert.Equal(t, "", passwordFlag.DefValue)

	usernameFlag := cmd.Flags().Lookup("username")
	assert.NotNil(t, usernameFlag)
	assert.Equal(t, "", usernameFlag.DefValue)
}
