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
	"github.com/spf13/cobra"
)

// UserDeleteCmd defines the "delete" command for user deletion.
func UserDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete user by name or id",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// If there are command line arguments, process them concurrently.
			if len(args) > 0 {
				var wg sync.WaitGroup
				errChan := make(chan error, len(args)) // Channel to collect errors

				for _, arg := range args {
					// Retrieve user ID by name.
					userID, err := api.GetUsersIdByName(arg)
					if err != nil {
						return fmt.Errorf("failed to get user id for '%s': %w", arg, err)
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
				var finalErr error
				for err := range errChan {
					if finalErr == nil {
						finalErr = err
					}
				}

				if finalErr != nil {
					if isUnauthorizedError(finalErr) {
						return fmt.Errorf("permission denied: admin privileges are required to execute this command")
					}
					return fmt.Errorf("failed to delete user: %w", finalErr)
				}
				return nil
			} else {
				// Interactive mode: get the user ID from the prompt.
				userID := prompt.GetUserIdFromUser()
				if err := api.DeleteUser(userID); err != nil {
					if isUnauthorizedError(err) {
						return fmt.Errorf("permission denied: admin privileges are required to execute this command")
					}
					return fmt.Errorf("failed to delete user: %w", err)
				}
				return nil
			}
		},
	}

	return cmd
}
