package ldap

import (
	"github.com/spf13/cobra"
)

func Ldap() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ldap",
		Short:   "Manage ldap users and groups",
		Example: `  harbor ldap ping`,
	}
	cmd.AddCommand(
		LdapSearchUserCmd(),
		LdapPingCmd(),
	)

	return cmd
}
