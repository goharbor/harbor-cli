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
