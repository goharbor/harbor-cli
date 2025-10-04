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
	"strings"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

// Import ldap users command
func LdapImportUserCmd() *cobra.Command {
	var uids []string

	cmd := &cobra.Command{
		Use:   "import [userID]",
		Short: "import ldap user by registered userid",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := api.LdapImportUser(append(uids, args[0]))
			if err != nil {
				return fmt.Errorf("failed to import ldap user: %v", err)
			}

			fmt.Println("Added users with ID: " + strings.Join(append(uids, args[0]), " "))
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&uids, "uid", "u", []string{}, "add more users to import")
	return cmd
}
