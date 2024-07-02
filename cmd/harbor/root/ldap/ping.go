package ldap

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Ping ldap server command
func LdapPingCmd() *cobra.Command {
	opts := &models.LdapConf{}
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "ping ldap server",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := api.LdapPingServer(opts)
			if err != nil {
				log.Fatalf("failed to ping ldap server: %v", err)
			}
			if response.Payload.Success {
				log.Info("Connection to LDAP Server Success")
			} else {
				log.Fatalf("connection to ldap server failed: %v", response.Payload.Message)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.LdapURL, "ldap-url", "", "", "URL of the ldap service")
	flags.StringVarP(&opts.LdapSearchPassword, "ldap-password", "", "", "search password of the ldap service")
	flags.StringVarP(&opts.LdapSearchDn, "ldap-search-dn", "", "", "User's dn who has the permission to search the ldap server")
	flags.StringVarP(&opts.LdapBaseDn, "ldap-base-dn", "", "", "The base dn from which to lookup the user")
	flags.StringVarP(&opts.LdapUID, "ldap-uid", "", "", "attribute used in search to match the user. It could be cn, uid based on your LDAP/AD.")
	flags.Int64VarP(&opts.LdapScope, "ldap-scope", "", 0, "search scope of ldap service default 0 base, 1 OneLevel, 2 Subtree.")
	flags.StringVarP(&opts.LdapFilter, "ldap-filter", "", "", "Search Filter of ldap service")
	flags.BoolVarP(&opts.LdapVerifyCert, "ldap-verify", "", false, "Verify Ldap server certificate")

	return cmd
}
