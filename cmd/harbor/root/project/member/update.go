package member

import (
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/spf13/cobra"
)

// NewGetRegistryCommand creates a new `harbor get registry` command
func UpdateMemberCommand() *cobra.Command {
	var opts api.UpdateMemberOptions
	var roleID int64

	cmd := &cobra.Command{
		Use:     "update [ProjectName Or ID] [member ID]",
		Short:   "update member by ID",
		Long:    "update member in a project by MemberID",
		Example: "  harbor project member update my-project [memberID] --roleid 2",
		Args:    cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) == 1 {
				opts.ProjectNameOrID = args[0]
			} else if len(args) == 2 {
				opts.ProjectNameOrID = args[0]
				opts.ID, _ = strconv.ParseInt(args[1], 0, 64)
			} else if opts.ProjectNameOrID == "" || opts.ID == 0 {
				if opts.ProjectNameOrID == "" {
					opts.ProjectNameOrID, err = prompt.GetProjectNameFromUser()
					if err != nil {
						return fmt.Errorf("failed to get project name: %v", err)
					}
				}
				if opts.ID == 0 {
					opts.ID = prompt.GetMemberIDFromUser(opts.ProjectNameOrID)
				}
			}

			if roleID == 0 {
				roleID = prompt.GetRoleIDFromUser()
			}
			opts.RoleID = &models.RoleRequest{
				RoleID: roleID,
			}

			err = api.UpdateMember(opts)
			if err != nil {
				return fmt.Errorf("failed to get members list: %v", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.ID, "id", "", 0, "Member ID")
	flags.Int64VarP(&roleID, "roleid", "", 0, "Role to be updated")
	flags.StringVarP(&opts.ProjectNameOrID, "projectname", "p", "", "Project Name")
	return cmd
}
