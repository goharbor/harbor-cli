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
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/user/create"
	"github.com/goharbor/harbor-cli/pkg/views/user/reset"

	log "github.com/sirupsen/logrus"
)

func CreateUser(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.User.CreateUser(ctx, &user.CreateUserParams{
		UserReq: &models.UserCreationReq{
			Email:    opts.Email,
			Realname: opts.Realname,
			Comment:  opts.Comment,
			Password: opts.Password,
			Username: opts.Username,
		},
	})
	if err != nil {
		return err
	}

	if response != nil {
		log.Infof("User `%s` created successfully", opts.Username)
	}

	return nil
}

func ResetPassword(userId int64, resetView reset.ResetView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.User.UpdateUserPassword(ctx, &user.UpdateUserPasswordParams{
		UserID: userId,
		Password: &models.PasswordReq{
			OldPassword: resetView.OldPassword,
			NewPassword: resetView.NewPassword,
		},
	})

	if err != nil {
		return err
	}
	log.Infof("User password reset successfully with id %d", userId)
	return nil
}

func DeleteUser(userId int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.User.DeleteUser(ctx, &user.DeleteUserParams{UserID: userId})
	if err != nil {
		return err
	}
	log.Infof("User deleted successfully with id %d", userId)
	return nil
}

func ElevateUser(userId int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	UserSysAdminFlag := &models.UserSysAdminFlag{
		SysadminFlag: true,
	}
	_, err = client.User.SetUserSysAdmin(ctx, &user.SetUserSysAdminParams{UserID: userId, SysadminFlag: UserSysAdminFlag})
	if err != nil {
		return err
	}
	log.Infof("user elevated role to admin successfully with id %d", userId)
	return nil
}

func ListUsers(opts ...ListFlags) (*user.ListUsersOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}
	var listFlags ListFlags
	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.User.ListUsers(ctx, &user.ListUsersParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Q:        &listFlags.Q,
		Sort:     &listFlags.Sort,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetUsersIdByName(userName string) (int64, error) {
	var opts ListFlags

	u, err := ListUsers(opts)
	if err != nil {
		return 0, err
	}
	for _, user := range u.Payload {
		if user.Username == userName {
			return user.UserID, nil
		}
	}

	return 0, fmt.Errorf("fail to get user Id by username: %s", userName)
}

func GetUserProfileById(userId int64) *models.UserResp {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil
	}

	resp, err := client.User.GetUser(ctx, &user.GetUserParams{UserID: userId})
	if err != nil {
		return nil
	}

	return resp.GetPayload()
}

func UpdateUserProfile(profile *models.UserProfile, userId int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.User.UpdateUserProfile(ctx, &user.UpdateUserProfileParams{
		UserID:  userId,
		Profile: profile,
	})

	if err != nil {
		return err
	}

	log.Infof("User's profile updated successfully with id %d", userId)

	return nil
}
