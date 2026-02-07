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
	"testing"

	"github.com/goharbor/harbor-cli/pkg/views/password/reset"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockUserPasswordChanger struct {
	id              map[string]int64
	passwords       map[int64]string
	userCnt         int
	expectAuthError bool
}

func (m *MockUserPasswordChanger) GetUserIDByName(username string) (int64, error) {
	if v, ok := m.id[username]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("username %s not found", username)
}

func (m *MockUserPasswordChanger) GetUserIDFromUser() int64 {
	return 999
}

func (m *MockUserPasswordChanger) FillPasswordView(resetView *reset.PasswordChangeView) {
	resetView.NewPassword = "NewPass456"
	resetView.ConfirmPassword = "NewPass456"
}

func (m *MockUserPasswordChanger) ResetPassword(userID int64, resetView reset.PasswordChangeView) error {
	if m.expectAuthError {
		return fmt.Errorf("403")
	}
	if _, ok := m.passwords[userID]; !ok {
		return fmt.Errorf("user %d not found", userID)
	}
	m.passwords[userID] = resetView.NewPassword
	return nil
}

func initMockUserPasswordChanger(userCnt int, expectAuthError bool) *MockUserPasswordChanger {
	m := &MockUserPasswordChanger{
		userCnt:         userCnt,
		expectAuthError: expectAuthError,
		id:              make(map[string]int64),
		passwords:       make(map[int64]string),
	}
	for i := 0; i < userCnt; i++ {
		m.id[fmt.Sprintf("test%d", i+1)] = int64(i + 1)
		m.passwords[int64(i+1)] = "InitialPass123"
	}
	return m
}

func TestChangePassword(t *testing.T) {
	tests := []struct {
		name                string
		setup               func() *MockUserPasswordChanger
		args                []string
		expectedPasswordID  int64
		expectedNewPassword string
		expectedErr         string
	}{
		{
			name: "successfully change password by username",
			setup: func() *MockUserPasswordChanger {
				return initMockUserPasswordChanger(5, false)
			},
			args:                []string{"test1"},
			expectedPasswordID:  1,
			expectedNewPassword: "NewPass456",
			expectedErr:         "",
		},
		{
			name: "change password via interactive prompt",
			setup: func() *MockUserPasswordChanger {
				m := initMockUserPasswordChanger(5, false)
				m.id["promptuser"] = 999
				m.passwords[999] = "InitialPass123"
				return m
			},
			args:                []string{},
			expectedPasswordID:  999,
			expectedNewPassword: "NewPass456",
			expectedErr:         "",
		},
		{
			name: "user not found logs error",
			setup: func() *MockUserPasswordChanger {
				return initMockUserPasswordChanger(5, false)
			},
			args:                []string{"nonexistent"},
			expectedPasswordID:  0,
			expectedNewPassword: "",
			expectedErr:         "failed to get user id",
		},
		{
			name: "user id is zero logs not found",
			setup: func() *MockUserPasswordChanger {
				m := initMockUserPasswordChanger(5, false)
				m.id["ghost"] = 0
				return m
			},
			args:                []string{"ghost"},
			expectedPasswordID:  0,
			expectedNewPassword: "",
			expectedErr:         "not found",
		},
		{
			name: "permission denied error",
			setup: func() *MockUserPasswordChanger {
				return initMockUserPasswordChanger(5, true)
			},
			args:                []string{"test1"},
			expectedPasswordID:  0,
			expectedNewPassword: "",
			expectedErr:         "Permission denied",
		},
		{
			name: "reset password fails with non-403 error",
			setup: func() *MockUserPasswordChanger {
				m := initMockUserPasswordChanger(5, false)
				delete(m.passwords, 1)
				return m
			},
			args:                []string{"test1"},
			expectedPasswordID:  0,
			expectedNewPassword: "",
			expectedErr:         "failed to reset user password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			originalLogOutput := log.StandardLogger().Out
			log.SetOutput(&buf)
			defer log.SetOutput(originalLogOutput)

			m := tt.setup()
			ChangePassword(tt.args, m)

			logs := buf.String()

			if tt.expectedErr != "" {
				assert.Contains(t, logs, tt.expectedErr, "Expected error logs to contain %s but got %s", tt.expectedErr, logs)
			} else {
				assert.Empty(t, logs, "Expected no error logs but got: %s", logs)
			}

			if tt.expectedPasswordID != 0 {
				password, exists := m.passwords[tt.expectedPasswordID]
				assert.True(t, exists, "User with ID %d should exist", tt.expectedPasswordID)
				assert.Equal(t, tt.expectedNewPassword, password, "Password for user %d should be changed to %s", tt.expectedPasswordID, tt.expectedNewPassword)
			}
		})
	}
}
