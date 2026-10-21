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
	"maps"
	"slices"
)

// Merge overlays configurations from left to right. Maps merge by key, named
// resources merge by name, and later explicitly specified fields take precedence.
func Merge(configurations ...*Configuration) (*Configuration, error) {
	result := NewConfiguration()
	for _, configuration := range configurations {
		if err := configuration.validate(false); err != nil {
			return nil, err
		}
		cloned, err := cloneConfiguration(configuration)
		if err != nil {
			return nil, err
		}
		mergeSpec(&result.Spec, cloned.Spec)
	}
	if err := result.Validate(); err != nil {
		return nil, err
	}
	result.Normalize()
	return result, nil
}

func cloneConfiguration(configuration *Configuration) (*Configuration, error) {
	encoded, err := json.Marshal(configuration)
	if err != nil {
		return nil, fmt.Errorf("clone configuration: %w", err)
	}
	var result Configuration
	if err := json.Unmarshal(encoded, &result); err != nil {
		return nil, fmt.Errorf("clone configuration: %w", err)
	}
	return &result, nil
}

func mergeSpec(destination *Spec, source Spec) {
	destination.System = mergeMap(destination.System, source.System)
	destination.Registries = mergeNamed(destination.Registries, source.Registries, func(value Registry) string { return value.Name }, mergeRegistry)
	destination.Projects = mergeNamed(destination.Projects, source.Projects, func(value Project) string { return value.Name }, mergeProject)
	destination.ReplicationPolicies = mergeNamed(
		destination.ReplicationPolicies,
		source.ReplicationPolicies,
		func(value ReplicationPolicy) string { return value.Name },
		mergeReplicationPolicy,
	)
}

func mergeNamed[T any](destination, source []T, name func(T) string, merge func(T, T) T) []T {
	indexes := make(map[string]int, len(destination))
	for i, value := range destination {
		indexes[name(value)] = i
	}
	for _, value := range source {
		index, exists := indexes[name(value)]
		if !exists {
			indexes[name(value)] = len(destination)
			destination = append(destination, value)
			continue
		}
		destination[index] = merge(destination[index], value)
	}
	return destination
}

func mergeRegistry(base, overlay Registry) Registry {
	base.Type = prefer(overlay.Type, base.Type)
	base.URL = prefer(overlay.URL, base.URL)
	base.Description = prefer(overlay.Description, base.Description)
	base.Insecure = prefer(overlay.Insecure, base.Insecure)
	if overlay.Credential != nil {
		if base.Credential == nil {
			base.Credential = &RegistryCredential{}
		}
		base.Credential.Type = prefer(overlay.Credential.Type, base.Credential.Type)
		base.Credential.AccessKeyFrom = prefer(overlay.Credential.AccessKeyFrom, base.Credential.AccessKeyFrom)
		base.Credential.AccessSecretFrom = prefer(overlay.Credential.AccessSecretFrom, base.Credential.AccessSecretFrom)
	}
	return base
}

func mergeProject(base, overlay Project) Project {
	base.Public = prefer(overlay.Public, base.Public)
	base.Registry = prefer(overlay.Registry, base.Registry)
	base.Metadata = mergeProjectMetadata(base.Metadata, overlay.Metadata)
	base.Quota = mergeMap(base.Quota, overlay.Quota)
	base.Webhooks = mergeNamed(base.Webhooks, overlay.Webhooks, func(value Webhook) string { return value.Name }, mergeWebhook)
	return base
}

func mergeProjectMetadata(base, overlay *ProjectMetadata) *ProjectMetadata {
	if overlay == nil {
		return base
	}
	if base == nil {
		base = &ProjectMetadata{}
	}
	base.AutoScan = prefer(overlay.AutoScan, base.AutoScan)
	base.AutoSBOMGeneration = prefer(overlay.AutoSBOMGeneration, base.AutoSBOMGeneration)
	base.EnableContentTrust = prefer(overlay.EnableContentTrust, base.EnableContentTrust)
	base.EnableContentTrustCosign = prefer(overlay.EnableContentTrustCosign, base.EnableContentTrustCosign)
	base.PreventVulnerableImages = prefer(overlay.PreventVulnerableImages, base.PreventVulnerableImages)
	base.ProxySpeedKB = prefer(overlay.ProxySpeedKB, base.ProxySpeedKB)
	base.ReuseSystemCVEAllowlist = prefer(overlay.ReuseSystemCVEAllowlist, base.ReuseSystemCVEAllowlist)
	base.Severity = prefer(overlay.Severity, base.Severity)
	return base
}

func mergeWebhook(base, overlay Webhook) Webhook {
	base.Description = prefer(overlay.Description, base.Description)
	base.Enabled = prefer(overlay.Enabled, base.Enabled)
	if overlay.Events != nil {
		base.Events = slices.Clone(overlay.Events)
	}
	if overlay.Targets != nil {
		base.Targets = slices.Clone(overlay.Targets)
	}
	return base
}

func mergeReplicationPolicy(base, overlay ReplicationPolicy) ReplicationPolicy {
	base.Description = prefer(overlay.Description, base.Description)
	base.Enabled = prefer(overlay.Enabled, base.Enabled)
	if overlay.Mode != "" {
		base.Mode = overlay.Mode
	}
	if overlay.Registry != "" {
		base.Registry = overlay.Registry
	}
	base.DestinationNamespace = prefer(overlay.DestinationNamespace, base.DestinationNamespace)
	base.DestinationReplaceCount = prefer(overlay.DestinationReplaceCount, base.DestinationReplaceCount)
	base.Override = prefer(overlay.Override, base.Override)
	base.ReplicateDeletion = prefer(overlay.ReplicateDeletion, base.ReplicateDeletion)
	base.CopyByChunk = prefer(overlay.CopyByChunk, base.CopyByChunk)
	base.Speed = prefer(overlay.Speed, base.Speed)
	if overlay.Filters != nil {
		base.Filters = slices.Clone(overlay.Filters)
	}
	if overlay.Trigger != nil {
		base.Trigger = overlay.Trigger
	}
	return base
}

func mergeMap[M ~map[K]V, K comparable, V any](destination, source M) M {
	if source == nil {
		return destination
	}
	if destination == nil {
		return maps.Clone(source)
	}
	maps.Copy(destination, source)
	return destination
}

func prefer[T any](preferred, fallback *T) *T {
	if preferred != nil {
		return preferred
	}
	return fallback
}
