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

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/replication"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/replication/policies/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func CreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create replication policies",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug("Starting replications create command")

			opts := &create.CreateView{}
			create.CreateRPolicyView(opts, false)
			registryID := prompt.GetRegistryNameFromUser()
			registry := api.GetRegistryResponse(registryID)
			policy := ConvertToPolicy(opts, registry)

			response, err := api.CreateReplicationPolicy(&replication.CreateReplicationPolicyParams{
				Policy: policy,
			})
			if err != nil {
				return fmt.Errorf("failed to create replication policy: %v", utils.ParseHarborErrorMsg(err))
			}
			fmt.Println("Replication policy created successfully with ID:", response.Location)
			return nil
		},
	}

	return cmd
}

func ConvertToPolicy(view *create.CreateView, registry *models.Registry) *models.ReplicationPolicy {
	policy := &models.ReplicationPolicy{
		Name:        view.Name,
		Description: view.Description,
		Enabled:     view.Enabled,
		Override:    view.Override,
		// ReplicateDeletion is the favored field to use for deletion replication
		// Deletion is deprecated and will be removed in future versions
		// However, for updating from false to true, we need to set both fields
		ReplicateDeletion: view.ReplicateDeletion,
		Deletion:          view.ReplicateDeletion,
		CopyByChunk:       &view.CopyByChunk,
	}

	if view.Speed != "" {
		speedInt, _ := strconv.ParseInt(view.Speed, 10, 32)
		speed := int32(speedInt)
		policy.Speed = &speed
	}

	trigger := &models.ReplicationTrigger{
		Type: view.TriggerType,
	}
	if view.TriggerType == "scheduled" {
		trigger.TriggerSettings = &models.ReplicationTriggerSettings{
			Cron: view.CronString,
		}
	}
	policy.Trigger = trigger

	if view.ReplicationMode == "Pull" {
		// Pull mode (external -> Harbor)
		policy.SrcRegistry = registry
		policy.DestRegistry = nil
	} else {
		// Push mode (Harbor -> external)
		policy.SrcRegistry = nil
		policy.DestRegistry = registry
	}

	return policy
}
