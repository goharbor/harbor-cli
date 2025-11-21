package usergroup

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"

	list "github.com/goharbor/harbor-cli/pkg/views/usergroup/list"
)

func UserGroupsListCommand() *cobra.Command {
	var (
		searchq  string
		searchID int64
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list user groups",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var userGroups []*models.UserGroup

			if searchID != -1 {
				resp, err := api.GetUserGroup(searchID)
				if err != nil {
					return err
				}

				userGroups = []*models.UserGroup{resp.Payload}
			} else if searchq != "" {
				_, err := api.SearchUserGroups(searchq)
				if err != nil {
					return err
				}

				// userGroups = []*models.UserGroup{&resp.Payload}
			} else {
				resp, err := api.ListUserGroups()
				if err != nil {
					return err
				}

				userGroups = resp.Payload
			}

			list.ListUserGroups(userGroups)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&searchq, "search", "-s", "", "use to search for a specific groupname")
	flags.Int64VarP(&searchID, "id", "-i", -1, "use to search for a specific groupid")

	return cmd
}
