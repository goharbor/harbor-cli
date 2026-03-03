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

var (
	getUserIDByName   = api.GetUsersIdByName
	getUserIDFromUser = prompt.GetUserIdFromUser
	fillPasswordView  = reset.ChangePasswordView
	resetPassword     = api.ResetPassword
)

func ChangePassword(args []string) {
	var userID int64
	var err error
	resetView := &reset.PasswordChangeView{}

	if len(args) > 0 {
		userID, err = getUserIDByName(args[0])
		if err != nil {
			log.Errorf("failed to get user id for '%s': %v", args[0], err)
			return
		}
		if userID == 0 {
			log.Errorf("user with name '%s' not found", args[0])
			return
		}
	} else {
		userID = getUserIDFromUser()
	}

	fillPasswordView(resetView)

	err = resetPassword(userID, *resetView)
	if err != nil {
		if isUnauthorizedError(err) {
			log.Error("permission denied: Admin privileges are required to execute this command.")
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
			ChangePassword(args)
		},
	}
	return cmd
}
