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

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete member by username",
		Long:    "delete members in a project by username of the member",
		Example: "  harbor project member delete my-project --username user",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) > 0 {
				ok, checkErr := api.CheckProject(args[0]) // verifying project name
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
				api.DeleteAllMember(project)
				fmt.Println("All members deleted successfully")
				return nil
			}

			if username != "" {
				err = api.DeleteMemberByUsername(project, username)
				if err != nil {
					return fmt.Errorf("failed to delete member: %v", err)
				}
				return nil
			} else if !delAllFlag {
				log.Println("Please provide a username or use --all flag to delete all members")
				memID = prompt.GetMemberIDFromUser(project)
				if memID == 0 {
					fmt.Println("No members found in project")
					return nil
				}
			}

			// normal deletion process
			err = api.DeleteMember(project, memID)
			if err != nil {
				return fmt.Errorf("failed to delete member: %v", err)
			}

			fmt.Printf("successfully deleted user with ID %d from project %s\n", memID, project)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&username, "username", "u", "", "Username of the member")
	flags.BoolVarP(&delAllFlag, "all", "a", false, "Deletes all members of the project")

	return cmd
}
