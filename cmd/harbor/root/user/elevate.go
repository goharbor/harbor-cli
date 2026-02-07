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

type UserElevator interface {
	GetUserIDByName(username string) (int64, error)
	GetUserIDFromUser() int64
	ConfirmElevation() (bool, error)
	ElevateUser(userID int64) error
}

type DefaultUserElevator struct{}

func (d *DefaultUserElevator) GetUserIDByName(username string) (int64, error) {
	return api.GetUsersIdByName(username)
}

func (d *DefaultUserElevator) GetUserIDFromUser() int64 {
	return prompt.GetUserIdFromUser()
}

func (d *DefaultUserElevator) ConfirmElevation() (bool, error) {
	return views.ConfirmElevation()
}

func (d *DefaultUserElevator) ElevateUser(userID int64) error {
	return api.ElevateUser(userID)
}

func ElevateUser(args []string, userElevator UserElevator) {
	var err error
	var userID int64

	if len(args) > 0 {
		userID, err = userElevator.GetUserIDByName(args[0])
		if err != nil {
			log.Errorf("failed to get user id for '%s': %v", args[0], err)
			return
		}
		if userID == 0 {
			log.Errorf("User with name '%s' not found", args[0])
			return
		}
	} else {
		userID = userElevator.GetUserIDFromUser()
	}

	confirm, err := userElevator.ConfirmElevation()
	if err != nil {
		log.Errorf("failed to confirm elevation: %v", err)
		return
	}
	if !confirm {
		log.Error("User did not confirm elevation. Aborting command.")
		return
	}

	err = userElevator.ElevateUser(userID)
	if err != nil {
		if isUnauthorizedError(err) {
			log.Error("Permission denied: Admin privileges are required to execute this command.")
		} else {
			log.Errorf("failed to elevate user: %v", err)
		}
	}
}

func ElevateUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "elevate",
		Short: "elevate user",
		Long:  "elevate user to admin role",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			d := &DefaultUserElevator{}
			ElevateUser(args, d)
		},
	}

	return cmd
}
