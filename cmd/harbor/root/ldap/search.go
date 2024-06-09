package ldap

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/ldap"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Search ldap users command
func LdapSearchUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [userID]",
		Short: "search ldap user by registered userid",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				err      error
				response *ldap.SearchLdapUserOK
			)

			response, err = api.LdapSearchUser(args[0])
			if err != nil {
				log.Fatalf("failed to search ldap user: %v", err)
			}

			utils.PrintPayloadInJSONFormat(response.Payload)
		},
	}

	return cmd
}
