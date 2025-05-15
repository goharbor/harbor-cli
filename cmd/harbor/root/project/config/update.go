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
package config

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/project/config/update"
	"github.com/spf13/cobra"
)

func UpdateProjectConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [project_name]",
		Short: "Update project configuration interactively",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectIDOrName string

			if len(args) > 0 {
				projectIDOrName = args[0]
			} else {
				projectIDOrName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("Failed to get project name: %v", err)
				}
				isID = false
			}

			resp, err := api.ListConfig(isID, projectIDOrName)
			if err != nil {
				return fmt.Errorf("Failed to list project config: %v", err)
			}
			config := resp.Payload
			conf := &models.ProjectMetadata{}
			if config != nil {
				for key, value := range config {
					switch key {
					case "public":
						conf.Public = value
					case "auto_scan":
						conf.AutoScan = &value
					case "prevent_vul":
						conf.PreventVul = &value
					case "reuse_sys_cve_allowlist":
						conf.ReuseSysCVEAllowlist = &value
					case "enable_content_trust":
						conf.EnableContentTrust = &value
					case "enable_content_trust_cosign":
						conf.EnableContentTrustCosign = &value
					case "severity":
						conf.Severity = &value
					}
				}
			}
			update.UpdateProjectMetadataView(conf)

			err = api.UpdateConfig(isID, projectIDOrName, *conf)
			if err != nil {
				return fmt.Errorf("Failed to update project config: %v", err)
			}
			return nil
		},
	}

	return cmd
}
