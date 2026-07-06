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
package policies

import (
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	config "github.com/goharbor/harbor-cli/pkg/config/replication"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/replication/policies/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// UpdateCommand returns a command to update existing replication policies
func UpdateCommand() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:   "update [policy-id]",
		Short: "Update an existing replication policy",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var policyID int64
			if len(args) > 0 {
				var err error
				policyID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return fmt.Errorf("invalid replication policy ID: %s, %v", args[0], err)
				}
			} else {
				policyID = prompt.GetReplicationPolicyFromUser()
			}

			var updatedPolicy *models.ReplicationPolicy

			if configFile != "" {
				log.Debugf("Loading replication policy configuration from file: %s", configFile)
				opts, err := config.LoadConfigFromFile(configFile)
				if err != nil {
					return fmt.Errorf("failed to load replication policy configuration: %v", err)
				}

				registryID, err := api.GetRegistryIdByName(opts.TargetRegistry)
				if err != nil {
					return fmt.Errorf("failed to get registry ID for name %s: %v", opts.TargetRegistry, err)
				}
				if registryID == 0 {
					return fmt.Errorf("registry with name %s not found", opts.TargetRegistry)
				}

				registry, err := api.GetRegistryResponse(registryID)
				if err != nil {
					return fmt.Errorf("failed to get registry with ID %d: %v", registryID, err)
				}

				updatedPolicy = ConvertToPolicy(opts, registry)
				updatedPolicy.ID = policyID
			} else {
				existingPolicy, err := api.GetReplicationPolicy(policyID)
				if err != nil {
					return fmt.Errorf("failed to get replication policy: %w", err)
				}

				var existingReplicationMode string
				if existingPolicy.Payload.SrcRegistry.ID != 0 && existingPolicy.Payload.DestRegistry.ID == 0 {
					existingReplicationMode = "Pull"
				} else if existingPolicy.Payload.SrcRegistry.ID == 0 && existingPolicy.Payload.DestRegistry.ID != 0 {
					existingReplicationMode = "Push"
				} else {
					return fmt.Errorf("replication policy with ID %d is neither Pull nor Push", policyID)
				}

				createView := &create.CreateView{
					Name:              existingPolicy.Payload.Name,
					Description:       existingPolicy.Payload.Description,
					Enabled:           existingPolicy.Payload.Enabled,
					Override:          existingPolicy.Payload.Override,
					ReplicateDeletion: existingPolicy.Payload.ReplicateDeletion,
					ReplicationMode:   existingReplicationMode,
				}

				if existingPolicy.Payload.CopyByChunk != nil {
					createView.CopyByChunk = *existingPolicy.Payload.CopyByChunk
				}

				if existingPolicy.Payload.Speed != nil {
					if *existingPolicy.Payload.Speed == 0 {
						speed := int32(-1)
						existingPolicy.Payload.Speed = &speed
					}
					createView.Speed = strconv.FormatInt(int64(*existingPolicy.Payload.Speed), 10)
				}

				if existingPolicy.Payload.SrcRegistry != nil && existingPolicy.Payload.DestRegistry == nil {
					createView.ReplicationMode = "Pull"
				} else if existingPolicy.Payload.SrcRegistry == nil && existingPolicy.Payload.DestRegistry != nil {
					createView.ReplicationMode = "Push"
				}

				if existingPolicy.Payload.Trigger != nil {
					createView.TriggerType = existingPolicy.Payload.Trigger.Type

					if existingPolicy.Payload.Trigger.TriggerSettings != nil {
						if existingPolicy.Payload.Trigger.Type == "scheduled" {
							createView.CronString = existingPolicy.Payload.Trigger.TriggerSettings.Cron
						} else if existingPolicy.Payload.Trigger.Type == "event_based" {
							createView.ReplicateDeletion = existingPolicy.Payload.ReplicateDeletion
						}
					}
				}

				log.Debugf("Updating replication policy: %s (ID: %d)", existingPolicy.Payload.Name, policyID)
				create.CreateRPolicyView(createView, true)

				fmt.Println("Updated policy replicate deletion:", createView.ReplicateDeletion)
				if createView.ReplicationMode == "Pull" {
					updatedPolicy = ConvertToPolicy(createView, existingPolicy.Payload.SrcRegistry)
					updatedPolicy.ID = policyID
				} else {
					updatedPolicy = ConvertToPolicy(createView, existingPolicy.Payload.DestRegistry)
				}
			}

			_, err := api.UpdateReplicationPolicy(policyID, updatedPolicy)
			if err != nil {
				return fmt.Errorf("failed to update replication policy: %v", utils.ParseHarborErrorMsg(err))
			}

			fmt.Printf("Successfully updated replication policy: %s (ID: %d)\n", updatedPolicy.Name, policyID)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&configFile, "policy-config-file", "f", "", "YAML/JSON file with replication policy configuration")

	return cmd
}
