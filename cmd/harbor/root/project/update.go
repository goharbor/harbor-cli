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
package project

import (
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/update"
	"github.com/spf13/cobra"
)

func UpdateProjectCommand() *cobra.Command {
	var (
		isID             bool
		publicFlag       string
		autoScanFlag     string
		preventVulFlag   string
		reuseSysCVEFlag  string
		severityFlag     string
		storageLimitFlag string
		registryIDFlag   string
	)

	cmd := &cobra.Command{
		Use:   "update [project_name]",
		Short: "Update a project",
		Long: `Update project settings such as visibility, storage quota, and metadata.

Examples:
  harbor project update myproject --public true
  harbor project update myproject --storage-limit -1 --prevent-vul true`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectNameOrID string

			if len(args) > 0 {
				projectNameOrID = args[0]
			} else {
				projectNameOrID, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", err)
				}
				isID = false
			}

			flags := cmd.Flags()
			flagsUsed := false
			storageChanged := false

			if flags.Changed("public") {
				if err := validateUpdateFlag("public", publicFlag); err != nil {
					return err
				}
				flagsUsed = true
			}
			if flags.Changed("auto-scan") {
				if err := validateUpdateFlag("auto-scan", autoScanFlag); err != nil {
					return err
				}
				flagsUsed = true
			}
			if flags.Changed("prevent-vul") {
				if err := validateUpdateFlag("prevent-vul", preventVulFlag); err != nil {
					return err
				}
				flagsUsed = true
			}
			if flags.Changed("reuse-sys-cve") {
				if err := validateUpdateFlag("reuse-sys-cve", reuseSysCVEFlag); err != nil {
					return err
				}
				flagsUsed = true
			}
			if flags.Changed("severity") {
				if err := validateUpdateFlag("severity", severityFlag); err != nil {
					return err
				}
				flagsUsed = true
			}
			if flags.Changed("storage-limit") {
				if err := utils.ValidateStorageLimit(storageLimitFlag); err != nil {
					return err
				}
				flagsUsed = true
				storageChanged = true
			}
			if flags.Changed("registry-id") {
				flagsUsed = true
			}

			resp, err := api.GetProject(projectNameOrID, isID)
			if err != nil {
				return fmt.Errorf("failed to get project: %v", utils.ParseHarborErrorMsg(err))
			}

			projectPayload := resp.Payload
			metadata := projectPayload.Metadata
			if metadata == nil {
				metadata = &models.ProjectMetadata{}
			}

			updateView := &update.UpdateView{
				Public:       metadata.Public,
				StorageLimit: currentStorageLimit(projectPayload.ProjectID),
				RegistryID:   strconv.FormatInt(projectPayload.RegistryID, 10),
				AutoScan:     metadata.AutoScan,
				PreventVul:   metadata.PreventVul,
				Severity:     metadata.Severity,
				ReuseSysCVE:  metadata.ReuseSysCVEAllowlist,
			}

			if flags.Changed("public") {
				metadata.Public = publicFlag
				updateView.Public = publicFlag
			}
			if flags.Changed("auto-scan") {
				metadata.AutoScan = &autoScanFlag
			}
			if flags.Changed("prevent-vul") {
				metadata.PreventVul = &preventVulFlag
			}
			if flags.Changed("reuse-sys-cve") {
				metadata.ReuseSysCVEAllowlist = &reuseSysCVEFlag
			}
			if flags.Changed("severity") {
				metadata.Severity = &severityFlag
			}
			if flags.Changed("storage-limit") {
				updateView.StorageLimit = storageLimitFlag
			}
			if flags.Changed("registry-id") {
				updateView.RegistryID = registryIDFlag
			}

			if !flagsUsed {
				if err := update.UpdateProjectView(updateView); err != nil {
					return fmt.Errorf("update cancelled or failed: %v", err)
				}
				metadata.Public = updateView.Public
				metadata.AutoScan = updateView.AutoScan
				metadata.PreventVul = updateView.PreventVul
				metadata.Severity = updateView.Severity
				metadata.ReuseSysCVEAllowlist = updateView.ReuseSysCVE
				storageChanged = true
			}

			public := metadata.Public == "true"
			req := &models.ProjectReq{
				Public:   &public,
				Metadata: metadata,
			}

			if updateView.RegistryID != "" && projectPayload.RegistryID != 0 {
				registryID, err := strconv.ParseInt(updateView.RegistryID, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid registry ID: %v", err)
				}
				req.RegistryID = &registryID
			}

			if err := api.UpdateProject(isID, projectNameOrID, req); err != nil {
				return fmt.Errorf("failed to update project: %v", utils.ParseHarborErrorMsg(err))
			}

			if storageChanged {
				storageLimit, err := strconv.ParseInt(updateView.StorageLimit, 10, 64)
				if err != nil {
					return fmt.Errorf("invalid storage limit: %v", err)
				}
				quota, err := api.GetQuotaByRef(int64(projectPayload.ProjectID))
				if err != nil {
					return fmt.Errorf("failed to get project quota: %v", utils.ParseHarborErrorMsg(err))
				}
				if err := api.UpdateQuota(quota.ID, &models.QuotaUpdateReq{
					Hard: models.ResourceList{"storage": storageLimit},
				}); err != nil {
					return fmt.Errorf("failed to update storage quota: %v", utils.ParseHarborErrorMsg(err))
				}
			}

			fmt.Printf("Project %s updated successfully\n", projectNameOrID)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&isID, "id", false, "Use project ID instead of name")
	flags.StringVar(&publicFlag, "public", "", "Set project visibility (true/false)")
	flags.StringVar(&autoScanFlag, "auto-scan", "", "Enable or disable auto scan (true/false)")
	flags.StringVar(&preventVulFlag, "prevent-vul", "", "Enable or disable vulnerability prevention (true/false)")
	flags.StringVar(&reuseSysCVEFlag, "reuse-sys-cve", "", "Enable or disable reuse of system CVE allowlist (true/false)")
	flags.StringVar(&severityFlag, "severity", "", "Set severity level (none, low, medium, high, critical)")
	flags.StringVar(&storageLimitFlag, "storage-limit", "", "Storage quota of the project in bytes (-1 for unlimited)")
	flags.StringVar(&registryIDFlag, "registry-id", "", "ID of referenced registry for proxy cache projects")

	return cmd
}

func currentStorageLimit(projectID int32) string {
	quota, err := api.GetQuotaByRef(int64(projectID))
	if err != nil {
		return "-1"
	}
	if val, ok := quota.Hard["storage"]; ok {
		return strconv.FormatInt(val, 10)
	}
	return "-1"
}

func validateUpdateFlag(flagName, flagValue string) error {
	allowed := map[string]bool{
		"none":     true,
		"low":      true,
		"medium":   true,
		"high":     true,
		"critical": true,
	}
	if flagName == "severity" && !allowed[flagValue] {
		return fmt.Errorf("invalid value for --%s: %s", flagName, flagValue)
	}
	if flagName != "severity" && flagValue != "true" && flagValue != "false" {
		return fmt.Errorf("invalid value for --%s: %s. Expected 'true' or 'false'", flagName, flagValue)
	}
	return nil
}
