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
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/replication/policies/create"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// updateOpts holds all non-interactive flag values for the update command.
type updateOpts struct {
	Name              string
	Description       string
	ResourceFilter    string
	NameFilter        string
	TagFilter         string
	TagPattern        string
	LabelFilter       string
	LabelPattern      string
	TriggerType       string
	CronString        string
	Speed             string
	Enabled           bool
	Override          bool
	ReplicateDeletion bool
	CopyByChunk       bool
}

// UpdateCommand returns a command to update existing replication policies
func UpdateCommand() *cobra.Command {
	var opts updateOpts

	cmd := &cobra.Command{
		Use:   "update [policy-id]",
		Short: "Update an existing replication policy",
		Long: `Update an existing replication policy.

When update flags are provided, the command runs non-interactively and updates only the specified fields while preserving all other values.`,
		Example: `harbor replication policies update 1 --name production-sync --enabled=true`,
		Args:    cobra.MaximumNArgs(1),
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
			if existingPolicy.Payload.SrcRegistry != nil && existingPolicy.Payload.SrcRegistry.ID != 0 &&
				(existingPolicy.Payload.DestRegistry == nil || existingPolicy.Payload.DestRegistry.ID == 0) {
				existingReplicationMode = "Pull"
			} else if (existingPolicy.Payload.SrcRegistry == nil || existingPolicy.Payload.SrcRegistry.ID == 0) &&
				existingPolicy.Payload.DestRegistry != nil && existingPolicy.Payload.DestRegistry.ID != 0 {
				existingReplicationMode = "Push"
			} else {
				return fmt.Errorf("replication policy with ID %d is neither Pull nor Push", policyID)
			}

			// Build the baseline CreateView from the existing policy.
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

			// Populate filter fields from the existing policy filters.
			for _, f := range existingPolicy.Payload.Filters {
				if f == nil {
					continue
				}
				switch f.Type {
				case "resource":
					if s, ok := f.Value.(string); ok {
						createView.ResourceFilter = s
					}
				case "name":
					if s, ok := f.Value.(string); ok {
						createView.NameFilter = s
					}
				case "tag":
					if s, ok := f.Value.(string); ok {
						createView.TagPattern = s
					}
					createView.TagFilter = f.Decoration
				case "label":
					createView.LabelFilter = f.Decoration
					switch v := f.Value.(type) {
					case string:
						createView.LabelPattern = v
					case []string:
						createView.LabelPattern = strings.Join(v, ",")
					case []interface{}:
						var parts []string
						for _, item := range v {
							if s, ok := item.(string); ok {
								parts = append(parts, s)
							}
						}
						createView.LabelPattern = strings.Join(parts, ",")
					}
				}
			}

			log.Debugf("Updating replication policy: %s (ID: %d)", existingPolicy.Payload.Name, policyID)

			// Branch: non-interactive if any update flag was explicitly provided.
			if hasReplicationUpdateFlagChanges(cmd) {
				if err := applyReplicationUpdateFlags(cmd, createView, opts); err != nil {
					return err
				}
			} else {
				// Preserve the existing interactive TUI workflow.
				create.CreateRPolicyView(createView, true)
			}

			var updatedPolicy *models.ReplicationPolicy
			if createView.ReplicationMode == "Pull" {
				updatedPolicy = ConvertToPolicy(createView, existingPolicy.Payload.SrcRegistry)
				updatedPolicy.ID = policyID
			} else {
				updatedPolicy = ConvertToPolicy(createView, existingPolicy.Payload.DestRegistry)
				updatedPolicy.ID = policyID
			}

			_, err = api.UpdateReplicationPolicy(policyID, updatedPolicy)
			if err != nil {
				return fmt.Errorf("failed to update replication policy: %w", err)
			}

			fmt.Printf("Successfully updated replication policy: %s (ID: %d)\n", updatedPolicy.Name, policyID)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Name, "name", "", "New name for the replication policy")
	flags.StringVar(&opts.Description, "description", "", "New description for the replication policy")
	flags.StringVar(&opts.ResourceFilter, "resource-filter", "", "Resource type filter: image, artifact, or empty for all")
	flags.StringVar(&opts.NameFilter, "name-filter", "", "Repository name filter pattern (supports wildcards, e.g. library/*)")
	flags.StringVar(&opts.TagFilter, "tag-filter", "", "Tag filter type: matches or excludes")
	flags.StringVar(&opts.TagPattern, "tag-pattern", "", "Tag filter pattern (e.g. v*, latest, *-prod)")
	flags.StringVar(&opts.LabelFilter, "label-filter", "", "Label filter type: matches or excludes")
	flags.StringVar(&opts.LabelPattern, "label-pattern", "", "Label filter pattern (e.g. env=prod or env=prod,ver=1.0)")
	flags.StringVar(&opts.TriggerType, "trigger-type", "", "Trigger type: manual, scheduled, or event_based")
	flags.StringVar(&opts.CronString, "cron", "", "Cron schedule (6-field format, required when --trigger-type=scheduled, e.g. \"0 0 */6 * * *\")")
	flags.StringVar(&opts.Speed, "speed", "", "Maximum replication speed in KB/s (-1 for unlimited)")
	flags.BoolVar(&opts.Enabled, "enabled", false, "Whether the replication policy is enabled or not")
	flags.BoolVar(&opts.Override, "override", false, "Override artifacts on destination if they already exist")
	flags.BoolVar(&opts.ReplicateDeletion, "replicate-deletion", false, "Replicate deletion operations to the destination")
	flags.BoolVar(&opts.CopyByChunk, "copy-by-chunk", false, "Transfer artifacts in chunks for better reliability")

	return cmd
}

// hasReplicationUpdateFlagChanges reports whether the user explicitly set any update flag.
func hasReplicationUpdateFlagChanges(cmd *cobra.Command) bool {
	flags := cmd.Flags()
	return flags.Changed("name") ||
		flags.Changed("description") ||
		flags.Changed("resource-filter") ||
		flags.Changed("name-filter") ||
		flags.Changed("tag-filter") ||
		flags.Changed("tag-pattern") ||
		flags.Changed("label-filter") ||
		flags.Changed("label-pattern") ||
		flags.Changed("trigger-type") ||
		flags.Changed("cron") ||
		flags.Changed("speed") ||
		flags.Changed("enabled") ||
		flags.Changed("override") ||
		flags.Changed("replicate-deletion") ||
		flags.Changed("copy-by-chunk")
}

// applyReplicationUpdateFlags overlays only the explicitly provided flags onto createView.
func applyReplicationUpdateFlags(cmd *cobra.Command, createView *create.CreateView, opts updateOpts) error {
	flags := cmd.Flags()

	if flags.Changed("name") {
		if strings.TrimSpace(opts.Name) == "" {
			return fmt.Errorf("--name cannot be empty")
		}
		createView.Name = strings.TrimSpace(opts.Name)
	}

	if flags.Changed("description") {
		createView.Description = opts.Description
	}

	if flags.Changed("resource-filter") {
		v := strings.ToLower(strings.TrimSpace(opts.ResourceFilter))
		switch createView.ReplicationMode {
		case "Pull":
			if v != "image" {
				return fmt.Errorf("--resource-filter must be 'image' for Pull mode, got %q", opts.ResourceFilter)
			}
		case "Push":
			if v != "" && v != "image" && v != "artifact" {
				return fmt.Errorf("--resource-filter must be '', 'image', or 'artifact' for Push mode, got %q", opts.ResourceFilter)
			}
		default:
			return fmt.Errorf("unknown replication mode %q", createView.ReplicationMode)
		}
		createView.ResourceFilter = v
	}

	if flags.Changed("name-filter") {
		createView.NameFilter = opts.NameFilter
	}

	if flags.Changed("tag-filter") {
		v := strings.ToLower(strings.TrimSpace(opts.TagFilter))
		if v != "matches" && v != "excludes" {
			return fmt.Errorf("--tag-filter must be 'matches' or 'excludes', got %q", opts.TagFilter)
		}
		createView.TagFilter = v
	}

	if flags.Changed("tag-pattern") {
		createView.TagPattern = opts.TagPattern
	}

	if flags.Changed("label-filter") {
		v := strings.ToLower(strings.TrimSpace(opts.LabelFilter))
		if v != "matches" && v != "excludes" {
			return fmt.Errorf("--label-filter must be 'matches' or 'excludes', got %q", opts.LabelFilter)
		}
		createView.LabelFilter = v
	}

	if flags.Changed("label-pattern") {
		createView.LabelPattern = opts.LabelPattern
	}

	if flags.Changed("trigger-type") {
		v := strings.ToLower(strings.TrimSpace(opts.TriggerType))
		if v != "manual" && v != "scheduled" && v != "event_based" {
			return fmt.Errorf("--trigger-type must be 'manual', 'scheduled', or 'event_based', got %q", opts.TriggerType)
		}
		createView.TriggerType = v
	}

	if flags.Changed("cron") {
		createView.CronString = strings.TrimSpace(opts.CronString)
		if createView.TriggerType != "scheduled" {
			return fmt.Errorf("--cron can only be used when --trigger-type=scheduled (current trigger-type: %s)", createView.TriggerType)
		}
	}

	// Validate cron dependency/format for scheduled triggers.
	if createView.TriggerType == "scheduled" {
		if createView.CronString == "" {
			return fmt.Errorf("--cron is required when --trigger-type=scheduled")
		}
		if flags.Changed("cron") || flags.Changed("trigger-type") {
			fields := strings.Fields(createView.CronString)
			if len(fields) != 6 {
				return fmt.Errorf("--cron must have exactly 6 fields (found %d): seconds minutes hours day-month month day-week", len(fields))
			}
		}
	}

	if flags.Changed("speed") {
		speedVal, err := strconv.ParseInt(opts.Speed, 10, 32)
		if err != nil || speedVal < -1 {
			return fmt.Errorf("--speed must be a valid integer >= -1 (use -1 for unlimited), got %q", opts.Speed)
		}
		createView.Speed = opts.Speed
	}

	if flags.Changed("enabled") {
		createView.Enabled = opts.Enabled
	}

	if flags.Changed("override") {
		createView.Override = opts.Override
	}

	if flags.Changed("replicate-deletion") {
		createView.ReplicateDeletion = opts.ReplicateDeletion
	}

	if flags.Changed("copy-by-chunk") {
		createView.CopyByChunk = opts.CopyByChunk
	}

	return nil
}
