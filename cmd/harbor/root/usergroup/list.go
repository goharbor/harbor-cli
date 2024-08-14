package usergroup

import (
    "github.com/goharbor/harbor-cli/pkg/api"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"

    list "github.com/goharbor/harbor-cli/pkg/views/usergroup/list"
)

func UserGroupsListCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "list",
        Short: "list user groups",
        Args:  cobra.NoArgs,
        Run: func(cmd *cobra.Command, args []string) {
            userGroups, err := api.ListUserGroups()
            if err != nil {
                log.Errorf("failed to list user groups: %v", err)
                return
            }

            list.ListUserGroups(userGroups)
        },
    }

    return cmd
}
