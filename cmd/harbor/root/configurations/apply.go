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
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	view "github.com/goharbor/harbor-cli/pkg/views/configurations/diff"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func ApplyConfigCmd() *cobra.Command {
	var cfgFile string
	var skipConfirm bool

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Update system configurations from local config file",
		Long: `Update Harbor system configurations using the values stored in your local config file.
		
This will push the configurations from your local config file to the Harbor server.
Make sure to run 'harbor config get' first to populate the local config file with current configurations. Alternatively, you can specify a custom configuration file using the --configurations-file flag. This does not have to be a complete configuration file, only the fields you want to update need to be present under the 'configurations' key. Credentials for the Harbor server can be configured in the local config file or through environment variables or global config flags.`,
		Args:    cobra.NoArgs,
		Example: `harbor config apply -f <config_file>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var configurations *models.Configurations
			var err error
			if cfgFile != "" {
				data, err := os.ReadFile(cfgFile)
				if err != nil {
					return fmt.Errorf("failed to read config file: %v", err)
				}
				fileType := filepath.Ext(cfgFile)
				switch fileType {
				case ".yaml", ".yml":
					if err := yaml.Unmarshal(data, &configurations); err != nil {
						return fmt.Errorf("failed to parse YAML: %v", err)
					}
				case ".json":
					if err := json.Unmarshal(data, &configurations); err != nil {
						return fmt.Errorf("failed to parse JSON: %v", err)
					}
				default:
					return fmt.Errorf("unsupported file type: %s, expected '.yaml/.yml' or '.json'", fileType)
				}
			} else {
				return fmt.Errorf("no config file specified")
			}

			response, err := api.GetConfigurations()
			if err != nil {
				return err
			}
			upstreamConfigs := utils.ExtractConfigValues(response.Payload) // *models.ConfigurationsResponse
			localConfigs := utils.ExtractConfigValues(configurations)      // *models.Configurations

			hasChanges := false
			for field, localVal := range localConfigs {
				upstreamVal, exists := upstreamConfigs[field]
				if !exists || fmt.Sprintf("%v", upstreamVal) != fmt.Sprintf("%v", localVal) {
					hasChanges = true
					break
				}
			}
			if !hasChanges {
				successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
				fmt.Println(successStyle.Render("✓ No changes detected."))
				return nil
			}
			// Show diff
			view.DiffConfigurations(upstreamConfigs, localConfigs)

			// Confirmation prompt
			if !skipConfirm {
				promptStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("3"))
				fmt.Print(promptStyle.Render("Do you want to apply these changes? (y/N): "))

				reader := bufio.NewReader(os.Stdin)
				userResponse, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("failed to read user input: %v", err)
				}

				userResponse = strings.TrimSpace(strings.ToLower(userResponse))
				if userResponse != "y" && userResponse != "yes" {
					cancelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
					fmt.Println(cancelStyle.Render("✗ Configuration update cancelled."))
					return nil
				}
			}

			err = api.UpdateConfigurations(configurations)
			if err != nil {
				return fmt.Errorf("failed to update Harbor configurations: %v", err)
			}

			fmt.Printf("harbor configurations updated successfully from %s.", cfgFile)
			return nil
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&cfgFile, "configurations-file", "f", "", "Harbor configurations file to apply.")

	return cmd
}
