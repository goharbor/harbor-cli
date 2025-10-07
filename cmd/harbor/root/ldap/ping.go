// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package ldap

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

// Ping ldap server command
func LdapPingCmd() *cobra.Command {
	opts := &models.LdapConf{}
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "ping ldap server",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := api.LdapPingServer(opts)
			if err != nil {
				return err
			}

			if response.Payload.Message != "" {
				fmt.Println(response.Payload.Message)
			}

			if response.Payload.Success {
				fmt.Println("Connection to LDAP Server Success")
			} else {
				return fmt.Errorf("connection to ldap server failed: %v", response.Payload.Message)
			}

			return nil
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
