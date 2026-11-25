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
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/password/reset"
	"github.com/goharbor/harbor-cli/pkg/views/user/create"
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
		fmt.Printf("User `%s` created successfully\n", opts.Username)
	}

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
	fmt.Printf("User deleted successfully with id %d\n", userId)
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
	fmt.Printf("user elevated role to admin successfully with id %d\n", userId)
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

// ErrUserNotFound is returned when a user name or ID matches no user.
var ErrUserNotFound = errors.New("not found")

// usernameQuery builds a q filter that matches userName exactly. Harbor
// URL-decodes q once more after the HTTP layer, so the name is escaped to
// keep characters such as '+' literal. Harbor also reads a value starting
// with [, { or ( as a range or list and drops one leading backslash from an
// exact value, so such names get an extra leading backslash.
func usernameQuery(userName string) string {
	if strings.ContainsAny(userName[:1], `\[{(`) {
		userName = `\` + userName
	}
	return "username=" + url.QueryEscape(userName)
}

// GetUsersIdByName returns the ID of the user whose name is exactly userName.
// It returns an error wrapping ErrUserNotFound when no user has that name.
func GetUsersIdByName(userName string) (int64, error) {
	if userName == "" {
		return 0, fmt.Errorf("user %q %w", userName, ErrUserNotFound)
	}

	const pageSize = 100
	q := usernameQuery(userName)
	for page := int64(1); ; page++ {
		u, err := ListUsers(ListFlags{Page: page, PageSize: pageSize, Q: q})
		if err != nil {
			return 0, err
		}
		for _, usr := range u.Payload {
			if usr.Username == userName {
				return usr.UserID, nil
			}
		}
		if len(u.Payload) < pageSize {
			return 0, fmt.Errorf("user %q %w", userName, ErrUserNotFound)
		}
	}
}

// userExists reports whether Harbor has a user with the given ID.
func userExists(id int64) (bool, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return false, err
	}
	_, err = client.User.GetUser(ctx, &user.GetUserParams{UserID: id})
	if err == nil {
		return true, nil
	}
	var notFound *user.GetUserNotFound
	if errors.As(err, &notFound) {
		return false, nil
	}
	return false, err
}

// ResolveUserID returns the ID of the user named nameOrID or, when nameOrID
// is a positive number, the user with that ID. A number that is the name of
// one user and the ID of another is rejected as ambiguous.
func ResolveUserID(nameOrID string) (int64, error) {
	nameID, err := GetUsersIdByName(nameOrID)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return 0, err
	}
	nameFound := err == nil

	id, parseErr := strconv.ParseInt(nameOrID, 10, 64)
	if parseErr != nil || id <= 0 || (nameFound && nameID == id) {
		return nameID, err
	}

	idFound, err := userExists(id)
	if err != nil {
		return 0, err
	}
	switch {
	case nameFound && idFound:
		return 0, fmt.Errorf("%q is ambiguous: it is the name of user ID %d and also a user ID", nameOrID, nameID)
	case nameFound:
		return nameID, nil
	case idFound:
		return id, nil
	}
	return 0, fmt.Errorf("user %q %w", nameOrID, ErrUserNotFound)
}

func ResetPassword(userId int64, opts reset.PasswordChangeView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.User.UpdateUserPassword(ctx, &user.UpdateUserPasswordParams{
		Password: &models.PasswordReq{
			NewPassword: opts.NewPassword,
		},
		UserID: userId,
	})
	if err != nil {
		return err
	}
	if response != nil {
		fmt.Printf("User password changed for userId %d\n", userId)
	}
	return nil
}
