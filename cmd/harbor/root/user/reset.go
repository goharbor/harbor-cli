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
	"github.com/goharbor/harbor-cli/pkg/views/user/reset"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UserResetCmd() *cobra.Command {
	var UserID int64
	cmd := &cobra.Command{
		Use:   "reset [username]",
		Short: "reset user's password",
		Long:  "Resets the password for a specific user by providing their username",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			resetView := &reset.ResetView{}
			flags := cmd.Flags()

			if len(args) > 0 {
				UserID, _ = api.GetUsersIdByName(args[0])
			} else {
				if !flags.Changed("userID") {
					UserID = prompt.GetUserIdFromUser()
				}
			}

			existingUser := api.GetUserProfileById(UserID)
			if existingUser == nil {
				logrus.Errorf("user is not found")
				return
			}

			reset.ResetUserView(resetView)

			err = api.ResetPassword(UserID, *resetView)

			if err != nil {
				logrus.Errorf("failed to reset user's password: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&UserID, "userID", "", -1, "ID of the user")

	return cmd
}
