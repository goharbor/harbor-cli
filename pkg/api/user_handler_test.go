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
package api

import (
	"context"
	"errors"
	"fmt"
	"testing"

	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserProfile(t *testing.T) {
	origContextWithClient := contextWithClientFunc
	defer func() { contextWithClientFunc = origContextWithClient }()
	contextWithClientFunc = func() (context.Context, *v2client.HarborAPI, error) {
		return nil, nil, errors.New("mocked error")
	}

	// Without a valid client context, this should return an error
	err := UpdateUserProfile(1, "test@example.com", "Test User", "Test Comment")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mocked error")
}

func TestGetUserByIDOrName(t *testing.T) {
	origContextWithClient := contextWithClientFunc
	defer func() { contextWithClientFunc = origContextWithClient }()
	contextWithClientFunc = func() (context.Context, *v2client.HarborAPI, error) {
		return nil, nil, errors.New("mocked error")
	}

	// Without a valid client context, this should return an error
	_, err := GetUserByIDOrName("admin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mocked error")
}

func TestGetUserByID(t *testing.T) {
	origContextWithClient := contextWithClientFunc
	defer func() { contextWithClientFunc = origContextWithClient }()
	contextWithClientFunc = func() (context.Context, *v2client.HarborAPI, error) {
		return nil, nil, errors.New("mocked error")
	}

	// Without a valid client context, this should return an error
	_, err := GetUserByID(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mocked error")
}

func TestGetUserByID_Success(t *testing.T) {
	origListUsersFunc := listUsersFunc
	defer func() { listUsersFunc = origListUsersFunc }()
	listUsersFunc = func(opts ...ListFlags) (*user.ListUsersOK, error) {
		return &user.ListUsersOK{
			Payload: []*models.UserResp{
				{UserID: 1, Username: "admin"},
			},
		}, nil
	}

	u, err := GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), u.UserID)
	assert.Equal(t, "admin", u.Username)
}

func TestGetUserByIDOrName_SuccessByName(t *testing.T) {
	origGetUsersIdByNameFunc := getUsersIdByNameFunc
	origGetUserByIDFunc := getUserByIDFunc
	defer func() {
		getUsersIdByNameFunc = origGetUsersIdByNameFunc
		getUserByIDFunc = origGetUserByIDFunc
	}()

	getUsersIdByNameFunc = func(userName string) (int64, error) {
		if userName == "testuser" {
			return 2, nil
		}
		return 0, errors.New("not found")
	}

	getUserByIDFunc = func(userID int64) (*models.UserResp, error) {
		if userID == 2 {
			return &models.UserResp{UserID: 2, Username: "testuser"}, nil
		}
		return nil, errors.New("not found")
	}

	u, err := GetUserByIDOrName("testuser")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), u.UserID)
	assert.Equal(t, "testuser", u.Username)
}

func TestGetUserByIDOrName_FallbackToNameWhenIDNotFound(t *testing.T) {
	origGetUsersIdByNameFunc := getUsersIdByNameFunc
	origGetUserByIDFunc := getUserByIDFunc
	defer func() {
		getUsersIdByNameFunc = origGetUsersIdByNameFunc
		getUserByIDFunc = origGetUserByIDFunc
	}()

	getUsersIdByNameFunc = func(userName string) (int64, error) {
		if userName == "2" {
			return 2, nil
		}
		return 0, errors.New("not found")
	}

	callCount := 0
	getUserByIDFunc = func(userID int64) (*models.UserResp, error) {
		callCount++
		// First call from strconv.ParseInt fails
		if callCount == 1 {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		// Second call after name lookup fallback succeeds
		if userID == 2 {
			return &models.UserResp{UserID: 2, Username: "2"}, nil
		}
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}

	u, err := GetUserByIDOrName("2")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), u.UserID)
	assert.Equal(t, "2", u.Username)
}
