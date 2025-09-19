package ldap

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// Search ldap users command
func LdapSearchUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [userID]",
		Short: "search ldap user by registered userid",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := api.LdapSearchUser(args[0])
			if err != nil {
				return fmt.Errorf("failed to search ldap user: %v", err)
			}

			utils.PrintPayloadInJSONFormat(response.Payload)
			return nil
		},
	}

	return cmd
}
