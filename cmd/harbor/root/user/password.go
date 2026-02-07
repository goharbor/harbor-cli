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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/password/reset"
	"github.com/spf13/cobra"
)

func UserPasswordChangeCmd() *cobra.Command {
	var opts reset.PasswordChangeView

	cmd := &cobra.Command{
		Use:   "password",
		Short: "Reset user password by name or id",
		Long:  "Allows admin to reset the password for a specified user or select interactively if no username is provided.",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			var userId int64
			var err error
			resetView := &reset.PasswordChangeView{
				NewPassword:     opts.NewPassword,
				ConfirmPassword: opts.ConfirmPassword,
			}

			if len(args) > 0 {
				userId, err = api.GetUsersIdByName(args[0])
				if err != nil {
					return fmt.Errorf("failed to get user id for '%s': %w", args[0], err)
				}
				if userId == 0 {
					return fmt.Errorf("user with name '%s' not found", args[0])
				}
			} else {
				userId = prompt.GetUserIdFromUser()
			}

			reset.ChangePasswordView(resetView)

			err = api.ResetPassword(userId, opts)
			if err != nil {
				if isUnauthorizedError(err) {
					return fmt.Errorf("permission denied: admin privileges are required to execute this command")
				}
				return fmt.Errorf("failed to reset user password: %w", err)
			}
			return nil
		},
	}
	return cmd
}
