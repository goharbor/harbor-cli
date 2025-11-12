package usergroup

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	delete "github.com/goharbor/harbor-cli/pkg/views/usergroup/delete"
	"github.com/spf13/cobra"
)

func UserGroupDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [groupID]",
		Short: "delete user group",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var groupID int64
			var err error

			if len(args) > 0 && args[0] != "" {
				groupID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to convert argument to ID: %v", err)
				}
			} else {
				response, err := api.ListUserGroups()
				if err != nil {
					return fmt.Errorf("failed to list user groups: %v", err)
				}

				opts, err := delete.DeleteUserGroupView(response)
				if err != nil {
					return fmt.Errorf("failed to obtain user selection: %v", err)
				}

				groupID = opts.ID
			}

			err = api.DeleteUserGroup(groupID)
			if err != nil {
				return fmt.Errorf("failed to delete user group: %v", err)
			}

			return nil
		},
	}

	return cmd
}
