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
	"net/mail"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/user/update"
	"github.com/spf13/cobra"
)

func UserUpdateCmd() *cobra.Command {
	var opts update.UpdateView

	cmd := &cobra.Command{
		Use:   "update [USER_NAME_OR_ID]",
		Short: "update user profile",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var userID int64
			var existingUser *models.UserResp

			if len(args) > 0 {
				existingUser, err = api.GetUserByIDOrName(args[0])
				if err != nil {
					return err
				}
				userID = existingUser.UserID
			} else {
				// Interactive mode: select user from list
				userID, err = prompt.GetUserIdFromUser()
				if err != nil {
					return fmt.Errorf("failed to get user id: %v", err)
				}
				existingUser, err = api.GetUserByID(userID)
				if err != nil {
					return err
				}
			}

			// If flags are provided, run non-interactively using the flags
			// If no flags are provided, open interactive UpdateView
			emailFlagSelected := cmd.Flags().Changed("email")
			realnameFlagSelected := cmd.Flags().Changed("realname")
			commentFlagSelected := cmd.Flags().Changed("comment")

			if emailFlagSelected || realnameFlagSelected || commentFlagSelected {
				// In non-interactive mode, use existing user values for flags not specified
				email := existingUser.Email
				if emailFlagSelected {
					email = opts.Email
				}
				realname := existingUser.Realname
				if realnameFlagSelected {
					realname = opts.Realname
				}
				comment := existingUser.Comment
				if commentFlagSelected {
					comment = opts.Comment
				}

				// Validate email format if it changed
				if email != "" && email != existingUser.Email {
					addr, err := mail.ParseAddress(email)
					if err != nil || addr.Address != email {
						return fmt.Errorf("invalid email format: %q", email)
					}
				}

				err = api.UpdateUserProfile(userID, email, realname, comment)
			} else {
				// Interactive mode
				updateView := &update.UpdateView{
					Email:    existingUser.Email,
					Realname: existingUser.Realname,
					Comment:  existingUser.Comment,
				}
				err = updateUserView(userID, updateView)
			}

			if err != nil {
				if isUnauthorizedError(err) {
					return fmt.Errorf("Permission denied: Admin privileges are required to execute this command.")
				} else {
					return fmt.Errorf("failed to update user: %v", err)
				}
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Email, "email", "", "", "Email")
	flags.StringVarP(&opts.Realname, "realname", "", "", "Realname")
	flags.StringVarP(&opts.Comment, "comment", "", "", "Comment")

	return cmd
}

func updateUserView(userID int64, updateView *update.UpdateView) error {
	update.UpdateUserView(updateView)
	return api.UpdateUserProfile(userID, updateView.Email, updateView.Realname, updateView.Comment)
}
