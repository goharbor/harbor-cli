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
package configurations

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/configurations/view"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ViewConfigCmd() *cobra.Command {
	var category string
	cmd := &cobra.Command{
		Use:   "view",
		Short: "View Harbor configurations",
		Long: `View Harbor system configurations. You can filter by category:
- authentication (auth): User and service authentication settings
- security (sec): Security policies and certificate settings  
- system (sys): General system behavior and storage settings

Examples:
  # View all configurations
  harbor config view
  
  # View only authentication configurations
  harbor config view --category authentication
  harbor config view --category auth
  
  # View only security configurations  
  harbor config view --category security
  harbor config view --category sec
  
  # View only system configurations
  harbor config view --category system
  harbor config view --category sys`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := api.GetConfigurations()
			if err != nil {
				return err
			}

			normalizedCategory := normalizeCategory(category)

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(response.Payload, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				view.ViewConfigurations(response.Payload, normalizedCategory)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter by category (authentication/auth, security/sec, system/sys)")

	return cmd
}

func normalizeCategory(category string) string {
	switch category {
	case "auth":
		return "authentication"
	case "sec":
		return "security"
	case "sys":
		return "system"
	default:
		return category
	}
}
