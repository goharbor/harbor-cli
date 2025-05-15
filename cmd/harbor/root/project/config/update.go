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

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/config/update"
	"github.com/spf13/cobra"
)

var (
	publicFlag                   string
	autoScanFlag                 string
	preventVulFlag               string
	reuseSysCVEAllowlistFlag     string
	enableContentTrustFlag       string
	enableContentTrustCosignFlag string
	severityFlag                 string
)

func UpdateProjectConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [project_name]",
		Short: "Interactively or via flags update project configuration in Harbor",
		Long: `Update the configuration settings of a Harbor project either interactively or directly using command-line flags.

You can specify the project by its name or ID as an argument. If not provided, you will be prompted to select a project interactively.

Examples:

  # Update project 'myproject' visibility to public
  harbor-cli project config update myproject --public true

  # Update multiple settings in one command
  harbor-cli project config update myproject --public false --prevent-vul true --severity high

  # Run interactively without flags
  harbor-cli project config update

Supported flag values:

  - Boolean flags (public, auto-scan, prevent-vul, reuse-sys-cve-allowlist, enable-content-trust, enable-content-trust-cosign): "true" or "false"
  - Severity: one of "low", "medium", "high", "critical"
`,

		Args: cobra.MaximumNArgs(1),
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
			resp, err := api.GetProject(projectIDOrName, isID)
			if err != nil {
				return fmt.Errorf("Failed to list project config: %v", utils.ParseHarborErrorMsg(err))
			}
			conf := resp.Payload.Metadata
			flags := cmd.Flags()
			flagsUsed := false

			if flags.Changed("public") {
				if err := validateFlag("public", publicFlag); err != nil {
					return err
				}
				conf.Public = publicFlag
				flagsUsed = true
			}
			if flags.Changed("auto-scan") {
				if err := validateFlag("auto-scan", autoScanFlag); err != nil {
					return err
				}
				conf.AutoScan = &autoScanFlag
				flagsUsed = true
			}
			if flags.Changed("prevent-vul") {
				if err := validateFlag("prevent-vul", preventVulFlag); err != nil {
					return err
				}
				conf.PreventVul = &preventVulFlag
				flagsUsed = true
			}
			if flags.Changed("reuse-sys-cve-allowlist") {
				if err := validateFlag("reuse-sys-cve-allowlist", reuseSysCVEAllowlistFlag); err != nil {
					return err
				}
				conf.ReuseSysCVEAllowlist = &reuseSysCVEAllowlistFlag
				flagsUsed = true
			}
			if flags.Changed("enable-content-trust") {
				if err := validateFlag("enable-content-trust", enableContentTrustFlag); err != nil {
					return err
				}
				conf.EnableContentTrust = &enableContentTrustFlag
				flagsUsed = true
			}
			if flags.Changed("enable-content-trust-cosign") {
				if err := validateFlag("enable-content-trust-cosign", enableContentTrustCosignFlag); err != nil {
					return err
				}
				conf.EnableContentTrustCosign = &enableContentTrustCosignFlag
				flagsUsed = true
			}
			if flags.Changed("severity") {
				if err := validateFlag("severity", severityFlag); err != nil {
					return err
				}
				conf.Severity = &severityFlag
				flagsUsed = true
			}
			if !flagsUsed {
				update.UpdateProjectMetadataView(conf)
			}

			err = api.UpdateConfig(isID, projectIDOrName, *conf)
			if err != nil {
				return fmt.Errorf("Failed to update project config: %v", utils.ParseHarborErrorMsg(err))
			}
			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVar(&publicFlag, "public", "", "Set project visibility (true/false)")
	flags.StringVar(&autoScanFlag, "auto-scan", "", "Enable or disable auto scan (true/false)")
	flags.StringVar(&preventVulFlag, "prevent-vul", "", "Enable or disable vulnerability prevention (true/false)")
	flags.StringVar(&reuseSysCVEAllowlistFlag, "reuse-sys-cve", "", "Enable or disable reuse of system CVE allowlist (true/false)")
	flags.StringVar(&enableContentTrustFlag, "enable-content-trust", "", "Enable or disable content trust (true/false)")
	flags.StringVar(&enableContentTrustCosignFlag, "enable-content-trust-cosign", "", "Enable or disable content trust cosign (true/false)")
	flags.StringVar(&severityFlag, "severity", "", "Set severity level")

	return cmd
}

func validateFlag(flagName, flagValue string) error {
	allowed := map[string]bool{
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}
	if flagName == "severity" && !allowed[flagValue] {
		return fmt.Errorf("Invalid value for --%s: %s. Allowed values are: low, medium, high, critical", flagName, flagValue)
	}
	if flagName != "severity" && flagValue != "true" && flagValue != "false" {
		return fmt.Errorf("Invalid value for --%s: %s. Expected 'true' or 'false'", flagName, flagValue)
	}

	return nil
}
