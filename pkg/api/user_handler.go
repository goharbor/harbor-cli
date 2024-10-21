package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/user/create"

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

func ListUsers() (*user.ListUsersOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.User.ListUsers(ctx, &user.ListUsersParams{})

	if err != nil {
		return nil, err
	}

	return response, nil
}
