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
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type mockUserDeleter struct {
	mu              sync.Mutex
	id              map[string]int64
	user            map[int64]string
	expectAuthError bool
}

func (m *mockUserDeleter) userDelete(userID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.expectAuthError {
		return fmt.Errorf("403")
	}
	if v, ok := m.user[userID]; ok {
		delete(m.id, v)
		delete(m.user, userID)
		return nil
	}
	return fmt.Errorf("user %d not found", userID)
}
func (m *mockUserDeleter) getUserIDByName(username string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.id[username]; ok {
		return v, nil
	} else {
		return 0, fmt.Errorf(`Username %s not found`, username)
	}
}
func (m *mockUserDeleter) getUserIDFromUser() int64 {
	return 999
}

func newMockUserDeleter(userCnt int, expectAuthError bool) *mockUserDeleter {
	m := &mockUserDeleter{
		expectAuthError: expectAuthError,
		id:              make(map[string]int64),
		user:            make(map[int64]string),
	}
	for i := 0; i < userCnt; i++ {
		m.id[fmt.Sprintf("test%d", i+1)] = int64(i + 1)
		m.user[int64(i+1)] = fmt.Sprintf("test%d", i+1)
	}
	getUserIDByName = m.getUserIDByName
	getUserIDFromUser = m.getUserIDFromUser
	userDeleteAPI = m.userDelete
	return m
}
func TestDeleteUser(t *testing.T) {
	origUserDelete := userDeleteAPI
	origGetUserIDByName := getUserIDByName
	origGetUserIDFromUser := getUserIDFromUser
	defer func() {
		userDeleteAPI = origUserDelete
		getUserIDByName = origGetUserIDByName
		getUserIDFromUser = origGetUserIDFromUser
	}()
	tests := []struct {
		name          string
		setup         func() *mockUserDeleter
		args          []string
		notExpectedID []int64
		expectedErr   string
	}{
		{
			name: "successfully delete single user",
			setup: func() *mockUserDeleter {
				return newMockUserDeleter(5, false)
			},
			args:          []string{"test1"},
			notExpectedID: []int64{1},
			expectedErr:   "",
		},
		{
			name: "successfully delete multiple users",
			setup: func() *mockUserDeleter {
				return newMockUserDeleter(5, false)
			},
			args:          []string{"test1", "test3", "test5"},
			notExpectedID: []int64{1, 3, 5},
			expectedErr:   "",
		},
		{
			name: "delete non-existent user logs error",
			setup: func() *mockUserDeleter {
				return newMockUserDeleter(5, false)
			},
			args:          []string{"nonexistent"},
			notExpectedID: []int64{},
			expectedErr:   "failed to get user id",
		},
		{
			name: "permission denied error",
			setup: func() *mockUserDeleter {
				return newMockUserDeleter(5, true)
			},
			args:          []string{"test1"},
			notExpectedID: []int64{},
			expectedErr:   "Permission denied",
		},
		{
			name: "mixed existing and non-existing users",
			setup: func() *mockUserDeleter {
				return newMockUserDeleter(5, false)
			},
			args:          []string{"test1", "nonexistent", "test3"},
			notExpectedID: []int64{1, 3},
			expectedErr:   "failed to get user id",
		},
		{
			name: "delete with empty args does not error",
			setup: func() *mockUserDeleter {
				m := newMockUserDeleter(5, false)
				m.user[999] = "promptuser"
				m.id["promptuser"] = 999
				return m
			},
			args:          []string{},
			notExpectedID: []int64{999},
			expectedErr:   "",
		},
		{
			name: "delete all users",
			setup: func() *mockUserDeleter {
				return newMockUserDeleter(3, false)
			},
			args:          []string{"test1", "test2", "test3"},
			notExpectedID: []int64{1, 2, 3},
			expectedErr:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			originalLogOutput := log.StandardLogger().Out
			log.SetOutput(&buf)
			defer log.SetOutput(originalLogOutput)

			m := tt.setup()
			DeleteUser(tt.args)

			logs := buf.String()

			if tt.expectedErr != "" {
				assert.Contains(t, logs, tt.expectedErr, "Expected error logs to contain %s but got %s", tt.expectedErr, logs)
			} else {
				assert.Empty(t, logs, "Expected no error logs but got: %s", logs)
			}

			for _, id := range tt.notExpectedID {
				_, ok := m.user[id]
				assert.False(t, ok, "User with ID %d should have been deleted", id)
			}
		})
	}
}
