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
	"github.com/goharbor/harbor-cli/pkg/views"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ElevateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elevate",
		Short: "elevate user",
		Long:  "elevate user to admin role",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var userId int64
			log.SetOutput(cmd.OutOrStderr())
			if len(args) > 0 {
				userId, err = api.GetUsersIdByName(args[0])
				if err != nil {
					log.Errorf("failed to get user id for '%s': %v", args[0], err)
					return
				}
				if userId == 0 {
					log.Errorf("User with name '%s' not found", args[0])
					return
				}
			} else {
				userId, err = prompt.GetUserIdFromUser()
				if err != nil {
					log.Errorf("failed to get user id: %v", err)
					return
				}
			}
			confirm, err := views.ConfirmElevation()
			if err != nil {
				log.Errorf("failed to confirm elevation: %v", err)
				return
			}
			if !confirm {
				log.Error("User did not confirm elevation. Aborting command.")
				return
			}

			err = api.ElevateUser(userId)
			if err != nil {
				if isUnauthorizedError(err) {
					log.Error("Permission denied: Admin privileges are required to execute this command.")
				} else {
					log.Errorf("failed to elevate user: %v", err)
				}
			}
		},
	}

	return cmd
}
