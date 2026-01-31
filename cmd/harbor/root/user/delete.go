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
	"sync"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type UserDeleter interface {
	UserDelete(userID int64) error
	GetUserIDByName(username string) (int64, error)
	GetUserIDFromUser() int64
}
type DefaultUserDeleter struct{}

func (d *DefaultUserDeleter) UserDelete(userID int64) error {
	return api.DeleteUser(userID)
}
func (d *DefaultUserDeleter) GetUserIDByName(username string) (int64, error) {
	return api.GetUsersIdByName(username)
}
func (d *DefaultUserDeleter) GetUserIDFromUser() int64 {
	return prompt.GetUserIdFromUser()
}

func DeleteUser(args []string, userDeleter UserDeleter) {
	// If there are command line arguments, process them concurrently.
	if len(args) > 0 {
		var wg sync.WaitGroup
		errChan := make(chan error, len(args)) // Channel to collect errors

		for _, arg := range args {
			// Retrieve user ID by name.
			userID, err := userDeleter.GetUserIDByName(arg)
			if err != nil {
				log.Errorf("failed to get user id for '%s': %v", arg, err)
				continue
			}
			wg.Add(1)
			go func(userID int64) {
				defer wg.Done()
				if err := userDeleter.UserDelete(userID); err != nil {
					errChan <- err
				}
			}(userID)
		}

		// Wait for all goroutines to finish and then close the error channel.
		go func() {
			wg.Wait()
			close(errChan)
		}()

		// Process errors from the goroutines.
		for err := range errChan {
			if isUnauthorizedError(err) {
				log.Error("Permission denied: Admin privileges are required to execute this command.")
			} else {
				log.Errorf("failed to delete user: %v", err)
			}
		}
	} else {
		// Interactive mode: get the user ID from the prompt.
		userID := userDeleter.GetUserIDFromUser()
		if err := userDeleter.UserDelete(userID); err != nil {
			if isUnauthorizedError(err) {
				log.Error("Permission denied: Admin privileges are required to execute this command.")
			} else {
				log.Errorf("failed to delete user: %v", err)
			}
		}
	}
}

// UserDeleteCmd defines the "delete" command for user deletion.
func UserDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete user by name or id", // nope it only deletes by name
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			d := &DefaultUserDeleter{}
			DeleteUser(args, d)
		},
	}

	return cmd
}
