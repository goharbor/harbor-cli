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
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/user/update"
	"github.com/spf13/cobra"
)

var (
	getUserByIDOrNameFunc = api.GetUserByIDOrName
	getUserIdFromUserFunc = prompt.GetUserIdFromUser
	getUserByIDFunc       = api.GetUserByID
	getCLIInfoFunc        = api.GetCLIInfo
	updateUserProfileFunc = api.UpdateUserProfile
	runUpdateUserViewFunc = update.UpdateUserView
)

func UserUpdateCmd() *cobra.Command {
	var opts update.UpdateView

	cmd := &cobra.Command{
		Use:     "update [USER_NAME_OR_ID]",
		Short:   "update user profile",
		Long:    `The 'update' command allows sysadmins to modify an existing user's profile information, such as their email, realname, or comment.`,
		Example: `  harbor user update admin --email newadmin@example.com --realname "System Admin"`,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var userID int64
			var existingUser *models.UserResp

			if len(args) > 0 {
				existingUser, err = getUserByIDOrNameFunc(args[0])
				if err != nil {
					return err
				}
				userID = existingUser.UserID
			} else {
				// Interactive mode: select user from list
				userID, err = getUserIdFromUserFunc()
				if err != nil {
					return fmt.Errorf("failed to get user id: %v", err)
				}
				existingUser, err = getUserByIDFunc(userID)
				if err != nil {
					return err
				}
			}

			cliInfo, err := getCLIInfoFunc()
			if err != nil {
				return err
			}

			if !cliInfo.IsSysAdmin && cliInfo.Username != existingUser.Username {
				return fmt.Errorf("permission denied: admin privileges are required to view or update other users")
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

				// Validate fields explicitly provided via flags
				if emailFlagSelected {
					addr, err := mail.ParseAddress(email)
					if err != nil || addr.Address != email {
						return fmt.Errorf("invalid email format: %q", email)
					}
				}
				if realnameFlagSelected && realname == "" {
					return fmt.Errorf("realname cannot be empty")
				}

				err = updateUserProfileFunc(userID, email, realname, comment)
			} else {
				// Interactive mode
				updateView := &update.UpdateView{
					Email:    existingUser.Email,
					Realname: existingUser.Realname,
					Comment:  existingUser.Comment,
				}
				err = runUpdateUserViewFunc(updateView)
				if err == nil {
					err = updateUserProfileFunc(userID, updateView.Email, updateView.Realname, updateView.Comment)
				}
			}

			if err != nil {
				return fmt.Errorf("failed to update user: %v", utils.ParseHarborErrorMsg(err))
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

