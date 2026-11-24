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
	"sync"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// UserDeleteCmd defines the "delete" command for user deletion.
func UserDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete user by name or id",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				var wg sync.WaitGroup
				errChan := make(chan error, len(args)) // Channel to collect errors

				// Fetch the user list once and resolve every name against it.
				resp, err := api.ListUsers()
				if err != nil {
					return fmt.Errorf("failed to list users: %w", err)
				}
				userIDs := make(map[string]int64, len(resp.Payload))
				for _, u := range resp.Payload {
					userIDs[u.Username] = u.UserID
				}

				for _, arg := range args {
					userID, ok := userIDs[arg]
					if !ok {
						log.Errorf("failed to get user id for '%s': user not found", arg)
						continue
					}
					wg.Add(1)
					go func(userID int64) {
						defer wg.Done()
						if err := api.DeleteUser(userID); err != nil {
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
						return fmt.Errorf("permission denied: admin privileges are required to execute this command")
					} else {
						return fmt.Errorf("failed to delete user: %w", err)
					}
				}
			} else {
				// Interactive mode: get the user ID from the prompt.
				userID, err := prompt.GetUserIdFromUser()
				if err != nil {
					return fmt.Errorf("failed to get user id: %v", err)
				}
				if err := api.DeleteUser(userID); err != nil {
					if isUnauthorizedError(err) {
						return fmt.Errorf("permission denied: Admin privileges are required to execute this command.")
					} else {
						return fmt.Errorf("failed to delete user: %v", err)
					}
				}
			}
			return nil
		},
	}

	return cmd
}
