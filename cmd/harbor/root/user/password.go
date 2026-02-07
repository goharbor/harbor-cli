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
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/password/reset"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type UserPasswordChanger interface {
	GetUserIDByName(username string) (int64, error)
	GetUserIDFromUser() int64
	FillPasswordView(resetView *reset.PasswordChangeView)
	ResetPassword(userID int64, resetView reset.PasswordChangeView) error
}

type DefaultUserPasswordChanger struct{}

func (d *DefaultUserPasswordChanger) GetUserIDByName(username string) (int64, error) {
	return api.GetUsersIdByName(username)
}

func (d *DefaultUserPasswordChanger) GetUserIDFromUser() int64 {
	return prompt.GetUserIdFromUser()
}

func (d *DefaultUserPasswordChanger) FillPasswordView(resetView *reset.PasswordChangeView) {
	reset.ChangePasswordView(resetView)
}

func (d *DefaultUserPasswordChanger) ResetPassword(userID int64, resetView reset.PasswordChangeView) error {
	return api.ResetPassword(userID, resetView)
}

func ChangePassword(args []string, passwordChanger UserPasswordChanger) {
	var userID int64
	var err error
	resetView := &reset.PasswordChangeView{}

	if len(args) > 0 {
		userID, err = passwordChanger.GetUserIDByName(args[0])
		if err != nil {
			log.Errorf("failed to get user id for '%s': %v", args[0], err)
			return
		}
		if userID == 0 {
			log.Errorf("User with name '%s' not found", args[0])
			return
		}
	} else {
		userID = passwordChanger.GetUserIDFromUser()
	}

	passwordChanger.FillPasswordView(resetView)

	err = passwordChanger.ResetPassword(userID, *resetView)
	if err != nil {
		if isUnauthorizedError(err) {
			log.Error("Permission denied: Admin privileges are required to execute this command.")
		} else {
			log.Errorf("failed to reset user password: %v", err)
		}
	}
}

func UserPasswordChangeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "password",
		Short: "Reset user password by name or id",
		Long:  "Allows admin to reset the password for a specified user or select interactively if no username is provided.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			d := &DefaultUserPasswordChanger{}
			ChangePassword(args, d)
		},
	}
	return cmd
}
