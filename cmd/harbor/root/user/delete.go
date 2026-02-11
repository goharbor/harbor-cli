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
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// UserDeleteCmd defines the "delete" command for user deletion.
func UserDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [username...]",
		Short: "delete user by name",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				userID := prompt.GetUserIdFromUser()
				if err := api.DeleteUser(userID); err != nil {
					return fmt.Errorf("failed to delete user: %v", utils.ParseHarborErrorMsg(err))
				}
				fmt.Println("User deleted successfully")
				return nil
			}

			allUsers, err := api.ListUsers()
			if err != nil {
				return fmt.Errorf("failed to list users: %v", utils.ParseHarborErrorMsg(err))
			}
			userMap := make(map[string]int64)
			for _, u := range allUsers.Payload {
				userMap[u.Username] = u.UserID
			}

			type user struct {
				name string
				id   int64
			}
			var resolved []user
			failedResolves := map[string]string{}

			for _, name := range args {
				id, found := userMap[name]
				if !found {
					failedResolves[name] = "user not found"
					continue
				}
				resolved = append(resolved, user{name: name, id: id})
			}

			var wg sync.WaitGroup
			var mu sync.Mutex
			successfulDeletes := []string{}
			failedDeletes := map[string]string{}

			for _, u := range resolved {
				wg.Add(1)
				go func(u user) {
					defer wg.Done()
					if err := api.DeleteUser(u.id); err != nil {
						mu.Lock()
						failedDeletes[u.name] = utils.ParseHarborErrorMsg(err)
						mu.Unlock()
					} else {
						mu.Lock()
						successfulDeletes = append(successfulDeletes, u.name)
						mu.Unlock()
					}
				}(u)
			}
			wg.Wait()

			if len(successfulDeletes) > 0 {
				fmt.Println("Successfully deleted users:")
				for _, name := range successfulDeletes {
					fmt.Printf("  - %s\n", name)
				}
			}

			if len(failedResolves) > 0 {
				fmt.Println("Failed to resolve users:")
				for name, reason := range failedResolves {
					fmt.Printf("  - %s: %s\n", name, reason)
				}
			}

			if len(failedDeletes) > 0 {
				fmt.Println("Failed to delete users:")
				for name, reason := range failedDeletes {
					fmt.Printf("  - %s: %s\n", name, reason)
				}
			}

			totFailed := len(failedResolves) + len(failedDeletes)
			if totFailed > 0 {
				return fmt.Errorf("failed to delete %d user(s)", totFailed)
			}

			log.Debug("All requested users processed for deletion successfully.")
			return nil
		},
	}

	return cmd
}
