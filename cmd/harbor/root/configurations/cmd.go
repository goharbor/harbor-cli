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
	"github.com/spf13/cobra"
)

func ConfigurationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage system configurations",
		Long: `Manage Harbor system configurations including viewing, exporting, and applying settings.

Configuration management workflow:
1. View configurations in table format or export to files
2. Edit exported configuration files as needed  
3. Apply modified configurations back to Harbor

Categories available:
- authentication (auth): LDAP, OIDC, UAA authentication settings
- security (sec): Security policies, certificates, and access control
- system (sys): General system behavior, storage, and operational settings`,
		Example: `  # View configurations
  harbor config view                          # Table view of all configs
  harbor config view -c auth                  # View only authentication configs
  harbor config view -c sec                   # View only security configs

  # Export configurations to files
  harbor config view -o json > config.json                    # Export all configs as JSON
  harbor config view -c auth -o yaml | tee auth-config.yaml   # Export auth configs as YAML
  harbor config view -c sys -o json > system-config.json     # Export system configs as JSON

  # Apply configurations from files
  harbor config apply -f config.json         # Apply complete configuration
  harbor config apply -f auth-config.yaml    # Apply only authentication settings
  
  # Configuration backup and restore workflow  
  harbor config view -o yaml > backup.yaml   # Create backup
  # ... make changes to Harbor via UI or other means ...
  harbor config apply -f backup.yaml         # Restore from backup`,
	}

	cmd.AddCommand(
		ViewConfigCmd(),
		ApplyConfigCmd(),
	)

	return cmd
}
