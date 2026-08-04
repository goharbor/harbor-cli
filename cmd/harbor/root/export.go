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

	"github.com/goharbor/harbor-cli/pkg/declarative"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func exportCommand() *cobra.Command {
	var filename string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export Harbor API-managed configuration",
		Long: `Export Harbor API-managed configuration as a versioned declarative document.

The export contains system settings, registry endpoints, projects, quotas,
webhooks, and replication policies. It does not contain credentials, robot
secrets, users, artifacts, database contents, or Harbor deployment settings.`,
		Example: `  harbor export -f harbor.yaml
  harbor export -f harbor.json -o json
  harbor export -o yaml > harbor.yaml`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			backend, err := declarative.NewLiveBackend()
			if err != nil {
				return err
			}
			configuration, err := declarative.NewService(backend).Export(cmd.Context())
			if err != nil {
				return err
			}

			format, err := exportFormat(filename, viper.GetString("output-format"))
			if err != nil {
				return err
			}
			if filename == "" || filename == "-" {
				return declarative.Encode(cmd.OutOrStdout(), configuration, format)
			}
			if err := declarative.WriteFile(filename, configuration, format); err != nil {
				return err
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Exported Harbor configuration to %s\n", filename)
			return nil
		},
	}
	cmd.Flags().StringVarP(&filename, "file", "f", "", "Write configuration to a YAML or JSON file (default: stdout)")
	return cmd
}

func exportFormat(filename, requested string) (declarative.Format, error) {
	if requested != "" {
		return declarative.ParseFormat(requested)
	}
	if filename != "" && filename != "-" {
		return declarative.FormatForFile(filename)
	}
	return declarative.FormatYAML, nil
}
