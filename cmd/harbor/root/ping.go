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
package root

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

// LoginCommand creates a new `harbor login` command

func PingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping the Harbor API server",
		Long:  "Send a ping request to the Harbor API server to check its status.",

		RunE: func(cmd *cobra.Command, args []string) error {
			err := api.Ping()
			if err != nil {
				return fmt.Errorf("failed to ping Harbor API server: %v", err)
			}
			fmt.Println("Harbor API server is reachable")
			return nil
		},
	}
	return cmd
}
