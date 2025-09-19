package ldap

import (
	"fmt"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

// Search ldap users command
func LdapImportUserCmd() *cobra.Command {
	var uids []string

	cmd := &cobra.Command{
		Use:   "import [userID]",
		Short: "import ldap user by registered userid",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := api.LdapImportUser(append(uids, args[0]))
			if err != nil {
				return fmt.Errorf("failed to search ldap user: %v", err)
			}

			fmt.Println("Added users with ID: " + strings.Join(append(uids, args[0]), " "))
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&uids, "uid", "u", []string{}, "add more users to import")
	return cmd
}
