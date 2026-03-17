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
	"fmt"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mockUserElevator struct {
	id                       map[string]int64
	admins                   map[int64]bool
	expectAuthError          bool
	confirmElevationFromUser bool
	confirmElevationErr      error
}

func (m *mockUserElevator) getUserIDByName(username string) (int64, error) {
	if v, ok := m.id[username]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("username %s not found", username)
}

func (m *mockUserElevator) getUserIDFromUser() (int64, error) {
	return 999, nil
}

func (m *mockUserElevator) confirmElevation() (bool, error) {
	return m.confirmElevationFromUser, m.confirmElevationErr
}

func (m *mockUserElevator) elevateUser(userID int64) error {
	if m.expectAuthError {
		return fmt.Errorf("403")
	}
	if _, ok := m.admins[userID]; !ok {
		m.admins[userID] = true
		return nil
	}
	return fmt.Errorf("user %d is already an admin", userID)
}

func newmockUserElevator(userCnt int, expectAuthError bool, confirmElevationFromUser bool, confirmElevationErr error) *mockUserElevator {
	m := &mockUserElevator{
		expectAuthError:          expectAuthError,
		confirmElevationFromUser: confirmElevationFromUser,
		confirmElevationErr:      confirmElevationErr,
		id:                       make(map[string]int64),
		admins:                   make(map[int64]bool),
	}
	for i := 0; i < userCnt; i++ {
		m.id[fmt.Sprintf("test%d", i+1)] = int64(i + 1)
	}
	getUserIDFromUser = m.getUserIDFromUser
	getUsersIDByName = m.getUserIDByName
	confirmElevation = m.confirmElevation
	elevateUserAPI = m.elevateUser
	return m
}

func TestElevateUser(t *testing.T) {
	origGetUsersId := getUsersIDByName
	origPrompt := getUserIDFromUser
	origConfirm := confirmElevation
	origElevate := elevateUserAPI
	defer func() {
		getUsersIDByName = origGetUsersId
		getUserIDFromUser = origPrompt
		confirmElevation = origConfirm
		elevateUserAPI = origElevate
	}()
	tests := []struct {
		name            string
		setup           func() *mockUserElevator
		args            []string
		expectedAdminID []int64
		expectedErr     string
	}{
		{
			name: "successfully elevate user by username",
			setup: func() *mockUserElevator {
				return newmockUserElevator(5, false, true, nil)
			},
			args:            []string{"test1"},
			expectedAdminID: []int64{1},
			expectedErr:     "",
		},
		{
			name: "elevate user via interactive prompt",
			setup: func() *mockUserElevator {
				m := newmockUserElevator(5, false, true, nil)
				m.id["promptuser"] = 999
				return m
			},
			args:            []string{},
			expectedAdminID: []int64{999},
			expectedErr:     "",
		},
		{
			name: "user not found logs error",
			setup: func() *mockUserElevator {
				return newmockUserElevator(5, false, true, nil)
			},
			args:            []string{"nonexistent"},
			expectedAdminID: []int64{},
			expectedErr:     "failed to get user id",
		},
		{
			name: "user id is zero logs not found",
			setup: func() *mockUserElevator {
				m := newmockUserElevator(5, false, true, nil)
				m.id["ghost"] = 0
				return m
			},
			args:            []string{"ghost"},
			expectedAdminID: []int64{},
			expectedErr:     "not found",
		},
		{
			name: "permission denied error",
			setup: func() *mockUserElevator {
				return newmockUserElevator(5, true, true, nil)
			},
			args:            []string{"test1"},
			expectedAdminID: []int64{},
			expectedErr:     "permission denied",
		},
		{
			name: "user declines elevation confirmation",
			setup: func() *mockUserElevator {
				return newmockUserElevator(5, false, false, nil)
			},
			args:            []string{"test1"},
			expectedAdminID: []int64{},
			expectedErr:     "did not confirm",
		},
		{
			name: "confirmation prompt returns error",
			setup: func() *mockUserElevator {
				return newmockUserElevator(5, false, false, fmt.Errorf("terminal error"))
			},
			args:            []string{"test1"},
			expectedAdminID: []int64{},
			expectedErr:     "failed to confirm",
		},
		{
			name: "elevate user that is already admin",
			setup: func() *mockUserElevator {
				m := newmockUserElevator(5, false, true, nil)
				m.admins[1] = true
				return m
			},
			args:            []string{"test1"},
			expectedAdminID: []int64{1},
			expectedErr:     "already an admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			originalLogOutput := log.StandardLogger().Out
			log.SetOutput(&buf)
			defer log.SetOutput(originalLogOutput)

			m := tt.setup()

			_ = ElevateUser(tt.args)

			logs := buf.String()
			logs = strings.ToLower(logs)

			if tt.expectedErr != "" {
				assert.Contains(t, logs, tt.expectedErr, "Expected error logs to contain %s but got %s", tt.expectedErr, logs)
			} else {
				assert.Empty(t, logs, "Expected no error logs but got: %s", logs)
			}

			for _, id := range tt.expectedAdminID {
				isAdmin, exists := m.admins[id]
				assert.True(t, exists && isAdmin, "User with ID %d should be an admin", id)
			}
		})
	}
}

func TestElevateUserCmd(t *testing.T) {
	cmd := ElevateUserCmd()

	assert.Equal(t, "elevate", cmd.Use)
	assert.Equal(t, "elevate user", cmd.Short)
	assert.Equal(t, "elevate user to admin role", cmd.Long)
	assert.NotNil(t, cmd.Args, "Args validator should be set")
}
