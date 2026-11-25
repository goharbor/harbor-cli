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
	"errors"
	"fmt"
	"sync"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

const deleteForbiddenMsg = "permission denied: admin privileges are required to execute this command"

// UserDeleteCmd defines the "delete" command for user deletion.
func UserDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [username|id]...",
		Short: "delete user by name or id",
		Long: `Delete one or more users by username or user ID, or pick one interactively when no argument is given.

Each argument is matched against usernames first. A number that is not a username is used as a user ID, and a number that is the username of one user and the ID of another is rejected as ambiguous. Every argument is resolved before any user is deleted, so nothing is deleted if one of them does not resolve.`,
		Args: cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				userID, err := prompt.GetUserIdFromUser()
				if err != nil {
					return fmt.Errorf("failed to get user id: %v", err)
				}
				return deleteUsers([]int64{userID})
			}

			// Resolve every argument before deleting anything.
			seen := make(map[int64]bool, len(args))
			userIDs := make([]int64, 0, len(args))
			for _, arg := range args {
				userID, err := api.ResolveUserID(arg)
				if err != nil {
					if isForbidden(err) {
						return errors.New(deleteForbiddenMsg)
					}
					return fmt.Errorf("failed to resolve user %q: %w", arg, err)
				}
				if !seen[userID] {
					seen[userID] = true
					userIDs = append(userIDs, userID)
				}
			}
			return deleteUsers(userIDs)
		},
	}

	return cmd
}

// deleteUsers deletes the users concurrently and waits for every deletion
// before reporting the first error.
func deleteUsers(userIDs []int64) error {
	var wg sync.WaitGroup
	errs := make([]error, len(userIDs))
	for i, userID := range userIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = api.DeleteUser(userID)
		}()
	}
	wg.Wait()

	for _, err := range errs {
		if err == nil {
			continue
		}
		if isForbidden(err) {
			return errors.New(deleteForbiddenMsg)
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// isForbidden reports whether err is a 403 response from the user API.
func isForbidden(err error) bool {
	var listErr *user.ListUsersForbidden
	var getErr *user.GetUserForbidden
	var deleteErr *user.DeleteUserForbidden
	return errors.As(err, &listErr) || errors.As(err, &getErr) || errors.As(err, &deleteErr)
}
