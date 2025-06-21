package policies

import (
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/replication/policies/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// UpdateCommand returns a command to update existing replication policies
func UpdateCommand() *cobra.Command {
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

			log.Infof("Updating replication policy: %s (ID: %d)", existingPolicy.Payload.Name, policyID)
			create.CreateRPolicyView(createView, true)

			var updatedPolicy *models.ReplicationPolicy
			if createView.ReplicationMode == "Pull" {
				updatedPolicy = ConvertToPolicy(createView, existingPolicy.Payload.SrcRegistry)
				updatedPolicy.ID = policyID
			} else {
				updatedPolicy = ConvertToPolicy(createView, existingPolicy.Payload.DestRegistry)
			}

			_, err = api.UpdateReplicationPolicy(policyID, updatedPolicy)
			if err != nil {
				return fmt.Errorf("failed to update replication policy: %w", err)
			}

			log.Infof("Successfully updated replication policy: %s (ID: %d)", updatedPolicy.Name, policyID)
			return nil
		},
	}

	return cmd
}
