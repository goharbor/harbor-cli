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
package member

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Deletes the member of the given project and Member
func DeleteMemberCommand() *cobra.Command {
	var delAllFlag bool
	var memID int64
	var project string
	var username string
	var isID bool

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete member by username",
		Long:    "delete members in a project by username of the member",
		Example: "  harbor project member delete my-project --username user",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) > 0 {
				ok, checkErr := api.CheckProject(args[0], isID) // verifying project name
				if checkErr != nil {
					return fmt.Errorf("failed to verify project name: %v", checkErr)
				}

				if ok {
					project = args[0]
				} else {
					return fmt.Errorf("invalid project name: %s", args[0])
				}
			} else {
				project, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
			}

			if delAllFlag {
				api.DeleteAllMember(project, !isID)
				fmt.Println("All members deleted successfully")
				return nil
			}

			if username != "" {
				err = api.DeleteMemberByUsername(project, username, !isID) // when set true parses projectNameOrID as projectName else it parses as an integer ID
				if err != nil {
					return fmt.Errorf("failed to delete member: %v", err)
				}
				return nil
			} else if !delAllFlag {
				log.Println("Please provide a username or use --all flag to delete all members")
				memID = prompt.GetMemberIDFromUser(project, username)
				if memID == 0 {
					fmt.Println("No members found in project")
					return nil
				}
			}

			// normal deletion process
			err = api.DeleteMember(project, memID, !isID) // when set true parses projectNameOrID as projectName else it parses as an integer ID
			if err != nil {
				return fmt.Errorf("failed to delete member: %v", err)
			}

			fmt.Printf("successfully deleted user with ID %d from project %s\n", memID, project)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "parses projectName as an ID")
	flags.StringVarP(&username, "username", "u", "", "Username of the member")
	flags.BoolVarP(&delAllFlag, "all", "a", false, "Deletes all members of the project")

	return cmd
}
