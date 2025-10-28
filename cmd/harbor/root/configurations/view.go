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
	"fmt"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/configure"
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
		Long: `View Harbor system configurations. You can filter by category using full names or shorthand:

Categories:
- authentication (auth): User and service authentication settings (LDAP, OIDC, UAA)
- security (sec): Security policies and certificate settings
- system (sys): General system behavior and storage settings

Examples:
  harbor config view                        # View all configurations
  harbor config view --category auth        # View authentication configs
  harbor config view --cat sec              # View security configs (shorthand)
  harbor config view --cat sys              # View system configs

  # Export configurations to files
  harbor config view -o json > config.json                    # Save all configs as JSON
  harbor config view --cat auth -o yaml | tee auth-config.yaml   # Save auth configs as YAML and display
  harbor config view --cat sec -o json > security-config.json   # Save security configs as JSON`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Expand shorthand category names
			expandedCategory := expandCategoryShorthand(category)

			// Validate category if provided
			if expandedCategory != "" {
				validCategories := []string{"authentication", "security", "system"}
				isValid := false
				for _, valid := range validCategories {
					if expandedCategory == valid {
						isValid = true
						break
					}
				}
				if !isValid {
					return fmt.Errorf("invalid category '%s'. Valid options: authentication (auth), security (sec), system (sys)", category)
				}
			}

			response, err := api.GetConfigurations()
			if err != nil {
				return err
			}

			var configurationsResponse *configure.GetConfigurationsOK

			configurations := utils.ExtractNonNullConfigurations(response.Payload)

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				configurations := utils.ExtractConfigurationsByCategory(configurations, expandedCategory)
				err = utils.PrintFormat(configurations, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				view.ViewConfigurations(configurations, expandedCategory)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter by category: authentication (auth), security (sec), system (sys)")
	cmd.Flags().StringVar(&category, "cat", "", "Filter by category (shorthand for --category)")

	return cmd
}

func expandCategoryShorthand(category string) string {
	switch strings.ToLower(category) {
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
