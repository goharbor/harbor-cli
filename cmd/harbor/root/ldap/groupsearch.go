package ldap

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Search ldap users command
func LdapSearchGroupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search-group [groupName or groupDN]",
		Short: "search ldap group by group name and group DN",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return fmt.Errorf("invalid groupName or groupDN provided")
			}

			response, err := api.LdapSearchGroup(args[0], "")
			if err != nil {
				return fmt.Errorf("failed to search ldap user: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(response.Payload, FormatFlag)
				if err != nil {
					return err
				}

				return nil
			}

			utils.PrintPayloadInJSONFormat(response.Payload)
			return nil
		},
	}

	return cmd
}
