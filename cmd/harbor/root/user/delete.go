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

func UserDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete user by name or id",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			var wg sync.WaitGroup
			errChan := make(chan error, len(args)) // Channel to collect error

			if len(args) > 0 {
				for _, arg := range args {
					userID, _ := api.GetUsersIdByName(arg)
					wg.Add(1)
					go func(userID int64) {
						defer wg.Done()
						if err := api.DeleteUser(userID); err != nil {
							errChan <- err
						}
					}(userID)
				}
			} else {
				userId := prompt.GetUserIdFromUser()
				err = api.DeleteUser(userId)
				if err != nil {
					log.Errorf("failed to delete user: %v", err)
				}
			}

			// Wait for all goroutines to finish
			go func() {
				wg.Wait()
				close(errChan)
			}()

			// Collect and handle errors
			var finalErr error
			for err := range errChan {
				if finalErr == nil {
					finalErr = err
				} else {
					log.Errorf("Error: %v", err)
				}
			}
			if finalErr != nil {
				log.Errorf("failed to delete user: %v", finalErr)
			}
		},
	}

	return cmd
}
