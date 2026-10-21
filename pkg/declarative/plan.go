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

package declarative

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
)

// Operation is the reconciliation operation for a resource.
type Operation string

const (
	// OperationCreate creates a missing resource.
	OperationCreate Operation = "create"
	// OperationUpdate updates a resource whose managed fields differ.
	OperationUpdate Operation = "update"
	// OperationNoop leaves an already-converged resource unchanged.
	OperationNoop Operation = "no-op"
)

const (
	resourceSystem            = "system configuration"
	resourceRegistry          = "registry"
	resourceProject           = "project"
	resourceQuota             = "quota"
	resourceWebhook           = "webhook"
	resourceReplicationPolicy = "replication policy"
)

// Action is one ordered reconciliation step.
type Action struct {
	Operation Operation
	Resource  string
	Name      string
	Project   string
}

// String returns a human-readable plan line.
func (a Action) String() string {
	identity := a.Name
	if a.Project != "" {
		identity = a.Project + "/" + a.Name
	}
	if identity == "" {
		return fmt.Sprintf("%s %s", a.Operation, a.Resource)
	}
	return fmt.Sprintf("%s %s %s", a.Operation, a.Resource, identity)
}

// Plan is an ordered set of reconciliation actions.
type Plan struct {
	Actions []Action
}

// HasChanges reports whether applying the plan would mutate Harbor.
func (p Plan) HasChanges() bool {
	return slices.ContainsFunc(p.Actions, func(action Action) bool { return action.Operation != OperationNoop })
}

// ChangeCount returns the number of mutating actions.
func (p Plan) ChangeCount() int {
	count := 0
	for _, action := range p.Actions {
		if action.Operation != OperationNoop {
			count++
		}
	}
	return count
}

// BuildPlan compares desired fields with current state and orders dependent actions.
func BuildPlan(desired, current *Configuration) (Plan, error) {
	if err := desired.Validate(); err != nil {
		return Plan{}, err
	}
	if current == nil {
		current = NewConfiguration()
	}

	registries := indexBy(current.Spec.Registries, func(value Registry) string { return value.Name })
	projects := indexBy(current.Spec.Projects, func(value Project) string { return value.Name })
	policies := indexBy(current.Spec.ReplicationPolicies, func(value ReplicationPolicy) string { return value.Name })

	var plan Plan
	for _, wanted := range desired.Spec.Registries {
		actual, exists := registries[wanted.Name]
		operation := OperationCreate
		if exists {
			operation = OperationNoop
			if !registryMatches(wanted, actual) {
				operation = OperationUpdate
			}
		}
		plan.Actions = append(plan.Actions, Action{Operation: operation, Resource: resourceRegistry, Name: wanted.Name})
	}

	for _, wanted := range desired.Spec.Projects {
		actual, exists := projects[wanted.Name]
		operation := OperationCreate
		if exists {
			operation = OperationNoop
			if !projectMatches(wanted, actual) {
				operation = OperationUpdate
			}
		}
		plan.Actions = append(plan.Actions, Action{Operation: operation, Resource: resourceProject, Name: wanted.Name})
	}

	for _, wanted := range desired.Spec.Projects {
		actual, projectExists := projects[wanted.Name]
		if wanted.Quota != nil {
			operation := OperationUpdate
			if projectExists && mapSubsetMatches(wanted.Quota, actual.Quota) {
				operation = OperationNoop
			}
			plan.Actions = append(plan.Actions, Action{Operation: operation, Resource: resourceQuota, Name: wanted.Name})
		}

		actualWebhooks := indexBy(actual.Webhooks, func(value Webhook) string { return value.Name })
		for _, webhook := range wanted.Webhooks {
			currentWebhook, exists := actualWebhooks[webhook.Name]
			operation := OperationCreate
			if projectExists && exists {
				operation = OperationNoop
				if !webhookMatches(webhook, currentWebhook) {
					operation = OperationUpdate
				}
			}
			plan.Actions = append(plan.Actions, Action{
				Operation: operation,
				Resource:  resourceWebhook,
				Name:      webhook.Name,
				Project:   wanted.Name,
			})
		}
	}

	for _, wanted := range desired.Spec.ReplicationPolicies {
		actual, exists := policies[wanted.Name]
		operation := OperationCreate
		if exists {
			operation = OperationNoop
			if !replicationPolicyMatches(wanted, actual) {
				operation = OperationUpdate
			}
		}
		plan.Actions = append(plan.Actions, Action{Operation: operation, Resource: resourceReplicationPolicy, Name: wanted.Name})
	}

	if len(desired.Spec.System) > 0 {
		operation := OperationNoop
		if !mapSubsetMatches(desired.Spec.System, current.Spec.System) {
			operation = OperationUpdate
		}
		plan.Actions = append(plan.Actions, Action{Operation: operation, Resource: resourceSystem})
	}
	return plan, nil
}

func indexBy[T any](values []T, key func(T) string) map[string]T {
	result := make(map[string]T, len(values))
	for _, value := range values {
		result[key(value)] = value
	}
	return result
}

func registryMatches(desired, current Registry) bool {
	return optionalMatches(desired.Type, current.Type) &&
		optionalMatches(desired.URL, current.URL) &&
		optionalMatches(desired.Description, current.Description) &&
		optionalMatches(desired.Insecure, current.Insecure) &&
		desired.Credential == nil
}

func projectMatches(desired, current Project) bool {
	return optionalMatches(desired.Public, current.Public) &&
		optionalMatches(desired.Registry, current.Registry) &&
		metadataMatches(desired.Metadata, current.Metadata)
}

func metadataMatches(desired, current *ProjectMetadata) bool {
	if desired == nil {
		return true
	}
	if current == nil {
		return false
	}
	return optionalMatches(desired.AutoScan, current.AutoScan) &&
		optionalMatches(desired.AutoSBOMGeneration, current.AutoSBOMGeneration) &&
		optionalMatches(desired.EnableContentTrust, current.EnableContentTrust) &&
		optionalMatches(desired.EnableContentTrustCosign, current.EnableContentTrustCosign) &&
		optionalMatches(desired.PreventVulnerableImages, current.PreventVulnerableImages) &&
		optionalMatches(desired.ProxySpeedKB, current.ProxySpeedKB) &&
		optionalMatches(desired.ReuseSystemCVEAllowlist, current.ReuseSystemCVEAllowlist) &&
		optionalMatches(desired.Severity, current.Severity)
}

func webhookMatches(desired, current Webhook) bool {
	return optionalMatches(desired.Description, current.Description) &&
		optionalMatches(desired.Enabled, current.Enabled) &&
		sliceMatches(desired.Events, current.Events) &&
		webhookTargetsMatch(desired.Targets, current.Targets)
}

func webhookTargetsMatch(desired, current []WebhookTarget) bool {
	if desired == nil {
		return true
	}
	if len(desired) != len(current) {
		return false
	}
	for i := range desired {
		wanted, actual := desired[i], current[i]
		if !optionalMatches(wanted.NotifyType, actual.NotifyType) ||
			!optionalMatches(wanted.Endpoint, actual.Endpoint) ||
			!optionalMatches(wanted.PayloadFormat, actual.PayloadFormat) ||
			!optionalMatches(wanted.VerifyRemoteCertificate, actual.VerifyRemoteCertificate) ||
			wanted.AuthHeaderFrom != nil {
			return false
		}
	}
	return true
}

func replicationPolicyMatches(desired, current ReplicationPolicy) bool {
	return desired.Mode == current.Mode && desired.Registry == current.Registry &&
		optionalMatches(desired.Description, current.Description) &&
		optionalMatches(desired.Enabled, current.Enabled) &&
		optionalMatches(desired.DestinationNamespace, current.DestinationNamespace) &&
		optionalMatches(desired.DestinationReplaceCount, current.DestinationReplaceCount) &&
		optionalMatches(desired.Override, current.Override) &&
		optionalMatches(desired.ReplicateDeletion, current.ReplicateDeletion) &&
		optionalMatches(desired.CopyByChunk, current.CopyByChunk) &&
		optionalMatches(desired.Speed, current.Speed) &&
		sliceMatches(desired.Filters, current.Filters) &&
		optionalMatches(desired.Trigger, current.Trigger)
}

func optionalMatches[T comparable](desired, current *T) bool {
	return desired == nil || current != nil && *desired == *current
}

func sliceMatches[T comparable](desired, current []T) bool {
	return desired == nil || slices.Equal(desired, current)
}

func mapSubsetMatches[K comparable, V any](desired, current map[K]V) bool {
	for key, wanted := range desired {
		actual, exists := current[key]
		if !exists || !jsonEquivalent(wanted, actual) {
			return false
		}
	}
	return true
}

func jsonEquivalent(left, right any) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	if leftErr != nil || rightErr != nil {
		return reflect.DeepEqual(left, right)
	}
	var normalizedLeft any
	var normalizedRight any
	if json.Unmarshal(leftJSON, &normalizedLeft) != nil || json.Unmarshal(rightJSON, &normalizedRight) != nil {
		return reflect.DeepEqual(left, right)
	}
	return reflect.DeepEqual(normalizedLeft, normalizedRight)
}
