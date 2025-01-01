package api

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/user/create"

	log "github.com/sirupsen/logrus"
)

func CreateUser(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client context for user: %s", opts.Username)
	}

	_, err = client.User.CreateUser(ctx, &user.CreateUserParams{
		UserReq: &models.UserCreationReq{
			Email:    opts.Email,
			Realname: opts.Realname,
			Comment:  opts.Comment,
			Password: opts.Password,
			Username: opts.Username,
		},
	})

	if err != nil {
		switch err.(type) {
		case *user.CreateUserBadRequest:
			return fmt.Errorf("bad request while creating user '%s'", opts.Username)
		case *user.CreateUserConflict:
			return fmt.Errorf("username or email already exists for user '%s'", opts.Username)
		case *user.CreateUserInternalServerError:
			return fmt.Errorf("internal server error occurred while creating user '%s'", opts.Username)
		default:
			return fmt.Errorf("an error occurred while creating user '%s'", opts.Username)
		}
	}

	log.Infof("User '%s' created successfully", opts.Username)
	return nil
}

func DeleteUser(userId int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client context for user ID: %d", userId)
	}

	_, err = client.User.DeleteUser(ctx, &user.DeleteUserParams{UserID: userId})
	if err != nil {
		switch err.(type) {
		case *user.DeleteUserNotFound:
			return fmt.Errorf("user not found: %d", userId)
		case *user.DeleteUserForbidden:
			return fmt.Errorf("insufficient permissions to delete user: %d", userId)
		case *user.DeleteUserUnauthorized:
			return fmt.Errorf("unauthorized access to delete user: %d", userId)
		case *user.DeleteUserInternalServerError:
			return fmt.Errorf("internal server error occurred while deleting user: %d", userId)
		default:
			return fmt.Errorf("an error occurred while deleting user: %d", userId)
		}
	}

	log.Infof("User with ID %d deleted successfully", userId)
	return nil
}

func ElevateUser(userId int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client context for user ID: %d", userId)
	}

	UserSysAdminFlag := &models.UserSysAdminFlag{
		SysadminFlag: true,
	}
	_, err = client.User.SetUserSysAdmin(ctx, &user.SetUserSysAdminParams{UserID: userId, SysadminFlag: UserSysAdminFlag})
	if err != nil {
		switch err.(type) {
		case *user.SetUserSysAdminNotFound:
			return fmt.Errorf("user not found: %d", userId)
		case *user.SetUserSysAdminInternalServerError:
			return fmt.Errorf("internal server error occurred while elevating user role to admin: %d", userId)
		case *user.SetUserSysAdminForbidden:
			return fmt.Errorf("insufficient permissions to elevate user role to admin: %d", userId)
		case *user.SetUserSysAdminUnauthorized:
			return fmt.Errorf("unauthorized access to elevate user role to admin: %d", userId)
		default:
			return fmt.Errorf("an error occurred while elevating user role to admin: %d", userId)
		}
	}

	log.Infof("User with ID %d elevated to admin successfully", userId)
	return nil
}

func ListUsers(opts ...ListFlags) (*user.ListUsersOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client context for listing users")
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
		switch err.(type) {
		case *user.ListUsersForbidden:
			return nil, fmt.Errorf("forbidden access while listing users")
		case *user.ListUsersUnauthorized:
			return nil, fmt.Errorf("unauthorized access while listing users")
		case *user.ListUsersInternalServerError:
			return nil, fmt.Errorf("internal server error occurred while listing users")
		default:
			return nil, fmt.Errorf("an error occurred while listing users")
		}
	}

	return response, nil
}

func GetUsersIdByName(userName string) (int64, error) {
	opts := ListFlags{}

	u, err := ListUsers(opts)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch user list for username: %s", userName)
	}

	for _, user := range u.Payload {
		if user.Username == userName {
			return user.UserID, nil
		}
	}

	return 0, fmt.Errorf("user not found: %s", userName)
}
