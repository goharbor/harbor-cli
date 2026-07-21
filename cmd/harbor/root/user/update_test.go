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
	"fmt"
	"testing"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/user/update"
	"github.com/stretchr/testify/assert"
)

func mockSetup(
	mockGetUserByIDOrName func(arg string) (*models.UserResp, error),
	mockGetUserIdFromUser func() (int64, error),
	mockGetUserByID func(userID int64) (*models.UserResp, error),
	mockGetCLIInfo func() (*api.CLIInfo, error),
	mockUpdateUserProfile func(userID int64, email, realname, comment string) error,
	mockRunUpdateUserView func(updateView *update.UpdateView) error,
) func() {
	origGetUserByIDOrName := getUserByIDOrNameFunc
	origGetUserIdFromUser := getUserIdFromUserFunc
	origGetUserByID := getUserByIDFunc
	origGetCLIInfo := getCLIInfoFunc
	origUpdateUserProfile := updateUserProfileFunc
	origRunUpdateUserView := runUpdateUserViewFunc

	getUserByIDOrNameFunc = mockGetUserByIDOrName
	getUserIdFromUserFunc = mockGetUserIdFromUser
	getUserByIDFunc = mockGetUserByID
	getCLIInfoFunc = mockGetCLIInfo
	updateUserProfileFunc = mockUpdateUserProfile
	runUpdateUserViewFunc = mockRunUpdateUserView

	return func() {
		getUserByIDOrNameFunc = origGetUserByIDOrName
		getUserIdFromUserFunc = origGetUserIdFromUser
		getUserByIDFunc = origGetUserByID
		getCLIInfoFunc = origGetCLIInfo
		updateUserProfileFunc = origUpdateUserProfile
		runUpdateUserViewFunc = origRunUpdateUserView
	}
}

func TestUserUpdateCmd_RunE(t *testing.T) {
	defaultUser := &models.UserResp{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
		Realname: "Test User",
		Comment:  "Comment",
	}
	defaultCLIInfo := &api.CLIInfo{
		IsSysAdmin: true,
		Username:   "admin",
	}

	tests := []struct {
		name        string
		args        []string
		flags       map[string]string
		setupMocks  func() func()
		expectedErr string
	}{
		{
			name: "non-interactive with valid flags",
			args: []string{"testuser"},
			flags: map[string]string{
				"email":    "new@example.com",
				"realname": "New Name",
				"comment":  "New Comment",
			},
			setupMocks: func() func() {
				return mockSetup(
					func(arg string) (*models.UserResp, error) { return defaultUser, nil },
					func() (int64, error) { return 0, nil },
					func(userID int64) (*models.UserResp, error) { return defaultUser, nil },
					func() (*api.CLIInfo, error) { return defaultCLIInfo, nil },
					func(userID int64, email, realname, comment string) error {
						if email != "new@example.com" {
							return fmt.Errorf("wrong email")
						}
						return nil
					},
					func(updateView *update.UpdateView) error { return nil },
				)
			},
			expectedErr: "",
		},
		{
			name: "non-interactive with realname and comment flags only",
			args: []string{"testuser"},
			flags: map[string]string{
				"realname": "New Name Only",
				"comment":  "New Comment Only",
			},
			setupMocks: func() func() {
				return mockSetup(
					func(arg string) (*models.UserResp, error) { return defaultUser, nil },
					func() (int64, error) { return 0, nil },
					func(userID int64) (*models.UserResp, error) { return defaultUser, nil },
					func() (*api.CLIInfo, error) { return defaultCLIInfo, nil },
					func(userID int64, email, realname, comment string) error {
						if realname != "New Name Only" {
							return fmt.Errorf("wrong realname")
						}
						if email != "test@example.com" { // should be existing user's email
							return fmt.Errorf("email should be unchanged")
						}
						return nil
					},
					func(updateView *update.UpdateView) error { return nil },
				)
			},
			expectedErr: "",
		},
		{
			name:  "interactive mode fallback",
			args:  []string{},
			flags: map[string]string{},
			setupMocks: func() func() {
				return mockSetup(
					func(arg string) (*models.UserResp, error) { return nil, nil },
					func() (int64, error) { return 1, nil },
					func(userID int64) (*models.UserResp, error) { return defaultUser, nil },
					func() (*api.CLIInfo, error) { return defaultCLIInfo, nil },
					func(userID int64, email, realname, comment string) error { return nil },
					func(updateView *update.UpdateView) error {
						updateView.Email = "interactive@example.com"
						return nil
					},
				)
			},
			expectedErr: "",
		},
		{
			name: "invalid email format",
			args: []string{"testuser"},
			flags: map[string]string{
				"email": "invalid-email",
			},
			setupMocks: func() func() {
				return mockSetup(
					func(arg string) (*models.UserResp, error) { return defaultUser, nil },
					func() (int64, error) { return 0, nil },
					func(userID int64) (*models.UserResp, error) { return defaultUser, nil },
					func() (*api.CLIInfo, error) { return defaultCLIInfo, nil },
					func(userID int64, email, realname, comment string) error { return nil },
					func(updateView *update.UpdateView) error { return nil },
				)
			},
			expectedErr: "invalid email format",
		},
		{
			name: "permission denied for non-admin",
			args: []string{"testuser"},
			flags: map[string]string{
				"email": "new@example.com",
			},
			setupMocks: func() func() {
				return mockSetup(
					func(arg string) (*models.UserResp, error) { return defaultUser, nil },
					func() (int64, error) { return 0, nil },
					func(userID int64) (*models.UserResp, error) { return defaultUser, nil },
					func() (*api.CLIInfo, error) {
						return &api.CLIInfo{IsSysAdmin: false, Username: "otheruser"}, nil
					},
					func(userID int64, email, realname, comment string) error { return nil },
					func(updateView *update.UpdateView) error { return nil },
				)
			},
			expectedErr: "permission denied",
		},
		{
			name: "self update for non-admin",
			args: []string{"testuser"},
			flags: map[string]string{
				"email": "new@example.com",
			},
			setupMocks: func() func() {
				return mockSetup(
					func(arg string) (*models.UserResp, error) { return defaultUser, nil },
					func() (int64, error) { return 0, nil },
					func(userID int64) (*models.UserResp, error) { return defaultUser, nil },
					func() (*api.CLIInfo, error) {
						return &api.CLIInfo{IsSysAdmin: false, Username: "testuser"}, nil
					},
					func(userID int64, email, realname, comment string) error { return nil },
					func(updateView *update.UpdateView) error { return nil },
				)
			},
			expectedErr: "",
		},
		{
			name:  "user resolution error",
			args:  []string{"testuser"},
			flags: map[string]string{},
			setupMocks: func() func() {
				return mockSetup(
					func(arg string) (*models.UserResp, error) { return nil, fmt.Errorf("user not found") },
					func() (int64, error) { return 0, nil },
					func(userID int64) (*models.UserResp, error) { return nil, nil },
					func() (*api.CLIInfo, error) { return defaultCLIInfo, nil },
					func(userID int64, email, realname, comment string) error { return nil },
					func(updateView *update.UpdateView) error { return nil },
				)
			},
			expectedErr: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := tt.setupMocks()
			defer cleanup()

			cmd := UserUpdateCmd()
			for k, v := range tt.flags {
				err := cmd.Flags().Set(k, v)
				assert.NoError(t, err)
			}

			err := cmd.RunE(cmd, tt.args)
			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserUpdateCmd_Metadata(t *testing.T) {
	cmd := UserUpdateCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "update [USER_NAME_OR_ID]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
	assert.NotEmpty(t, cmd.Example)
	assert.NotNil(t, cmd.RunE)
}

func TestUserUpdateCmd_Flags(t *testing.T) {
	cmd := UserUpdateCmd()
	assert.NotNil(t, cmd.Flags().Lookup("email"))
	assert.NotNil(t, cmd.Flags().Lookup("realname"))
	assert.NotNil(t, cmd.Flags().Lookup("comment"))
}
