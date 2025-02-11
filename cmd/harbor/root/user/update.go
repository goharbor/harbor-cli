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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/user/update"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateUserCmd() *cobra.Command {
	var UserID int64
	opts := &models.UserProfile{}

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "update user's profile",
		Long:    "Update a user's profile by providing their username, allowing you to modify personal details such as realname, email, or comment",
		Example: "harbor user update [username]",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			flags := cmd.Flags()

			if len(args) > 0 {
				UserID, err = api.GetUsersIdByName(args[0])
				if err != nil {
					logrus.Errorf("fail to get user id by username: %s", args[0])
					return
				}
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

			updateView := &models.UserProfile{
				Comment:  existingUser.Comment,
				Email:    existingUser.Email,
				Realname: existingUser.Realname,
			}

			if flags.Changed("comment") {
				updateView.Comment = opts.Comment
			}
			if flags.Changed("email") {
				updateView.Email = opts.Email
			}
			if flags.Changed("realname") {
				updateView.Realname = opts.Realname
			}

			update.UpdateUserProfileView(updateView)
			err = api.UpdateUserProfile(updateView, UserID)
			if err != nil {
				logrus.Errorf("failed to update user's profile: %v", err)
			}
		},
	}
	flags := cmd.Flags()
	flags.Int64VarP(&UserID, "userID", "", -1, "ID of the user")
	flags.StringVarP(&opts.Comment, "comment", "", "", "Comment of the user")
	flags.StringVarP(&opts.Email, "email", "", "", "Email of the user")
	flags.StringVarP(&opts.Realname, "realname", "", "", "Realname of the user")

	return cmd
}
