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
		Example: "  harbor project member delete --project my-project --username user",
		Args:    cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if project == "" {
				project, err = prompt.GetProjectNameFromUser()
				if err != nil {
					log.Fatalf("failed to get project name: %v", err)
				}
			}

			if delAllFlag {
				api.DeleteAllMember(project)
				fmt.Println("All members deleted successfully")
				return
			}

			if username != "" {
				err = api.DeleteMemberByUsername(project, username)
				if err != nil {
					log.Fatalf("failed to delete member: %v", err)
				}
				return
			} else if !delAllFlag {
				log.Println("Please provide a username or use --all flag to delete all members")
				memID = prompt.GetMemberIDFromUser(project)
				if memID == 0 {
					fmt.Println("No members found in project")
					return
				}
			}

			// normal deletion process
			err = api.DeleteMember(project, memID)
			if err != nil {
				log.Fatalf("failed to delete member: %v", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&username, "username", "u", "", "Username of the member")
	flags.StringVarP(&project, "project", "p", "", "Project Name")
	flags.BoolVarP(&delAllFlag, "all", "a", false, "Deletes all members of the project")

	return cmd
}
