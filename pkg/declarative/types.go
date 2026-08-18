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

// Package declarative exports and reconciles Harbor API-managed configuration.
package declarative

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

const (
	// APIVersion is the version of the declarative configuration schema.
	APIVersion = "goharbor.io/v1alpha1"
	// Kind identifies a complete Harbor configuration document.
	Kind = "HarborConfiguration"
)

// Configuration is the portable desired state of Harbor API-managed resources.
type Configuration struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	Spec       Spec   `json:"spec" yaml:"spec"`
}

// Spec contains the resources managed by a Configuration.
type Spec struct {
	System              map[string]any      `json:"system,omitempty" yaml:"system,omitempty"`
	Registries          []Registry          `json:"registries,omitempty" yaml:"registries,omitempty"`
	Projects            []Project           `json:"projects,omitempty" yaml:"projects,omitempty"`
	ReplicationPolicies []ReplicationPolicy `json:"replicationPolicies,omitempty" yaml:"replicationPolicies,omitempty"`
}

// SecretRef resolves a secret from the process environment at apply time.
type SecretRef struct {
	Env string `json:"env" yaml:"env"`
}

// Registry describes an external registry endpoint. Export never includes credentials.
type Registry struct {
	Name        string              `json:"name" yaml:"name"`
	Type        *string             `json:"type,omitempty" yaml:"type,omitempty"`
	URL         *string             `json:"url,omitempty" yaml:"url,omitempty"`
	Description *string             `json:"description,omitempty" yaml:"description,omitempty"`
	Insecure    *bool               `json:"insecure,omitempty" yaml:"insecure,omitempty"`
	Credential  *RegistryCredential `json:"credential,omitempty" yaml:"credential,omitempty"`
}

// RegistryCredential references external registry credentials without storing secret values.
type RegistryCredential struct {
	Type             *string    `json:"type,omitempty" yaml:"type,omitempty"`
	AccessKeyFrom    *SecretRef `json:"accessKeyFrom,omitempty" yaml:"accessKeyFrom,omitempty"`
	AccessSecretFrom *SecretRef `json:"accessSecretFrom,omitempty" yaml:"accessSecretFrom,omitempty"`
}

// Project describes a Harbor project and its project-scoped configuration.
type Project struct {
	Name     string           `json:"name" yaml:"name"`
	Public   *bool            `json:"public,omitempty" yaml:"public,omitempty"`
	Registry *string          `json:"registry,omitempty" yaml:"registry,omitempty"`
	Metadata *ProjectMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Quota    map[string]int64 `json:"quota,omitempty" yaml:"quota,omitempty"`
	Webhooks []Webhook        `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`
}

// ProjectMetadata contains mutable project settings.
type ProjectMetadata struct {
	AutoScan                 *bool   `json:"autoScan,omitempty" yaml:"autoScan,omitempty"`
	AutoSBOMGeneration       *bool   `json:"autoSBOMGeneration,omitempty" yaml:"autoSBOMGeneration,omitempty"`
	EnableContentTrust       *bool   `json:"enableContentTrust,omitempty" yaml:"enableContentTrust,omitempty"`
	EnableContentTrustCosign *bool   `json:"enableContentTrustCosign,omitempty" yaml:"enableContentTrustCosign,omitempty"`
	PreventVulnerableImages  *bool   `json:"preventVulnerableImages,omitempty" yaml:"preventVulnerableImages,omitempty"`
	ProxySpeedKB             *int64  `json:"proxySpeedKB,omitempty" yaml:"proxySpeedKB,omitempty"`
	ReuseSystemCVEAllowlist  *bool   `json:"reuseSystemCVEAllowlist,omitempty" yaml:"reuseSystemCVEAllowlist,omitempty"`
	Severity                 *string `json:"severity,omitempty" yaml:"severity,omitempty"`
}

// Webhook describes a project webhook policy.
type Webhook struct {
	Name        string          `json:"name" yaml:"name"`
	Description *string         `json:"description,omitempty" yaml:"description,omitempty"`
	Enabled     *bool           `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Events      []string        `json:"events,omitempty" yaml:"events,omitempty"`
	Targets     []WebhookTarget `json:"targets,omitempty" yaml:"targets,omitempty"`
}

// WebhookTarget describes one destination of a webhook policy.
type WebhookTarget struct {
	NotifyType              *string    `json:"notifyType,omitempty" yaml:"notifyType,omitempty"`
	Endpoint                *string    `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	PayloadFormat           *string    `json:"payloadFormat,omitempty" yaml:"payloadFormat,omitempty"`
	VerifyRemoteCertificate *bool      `json:"verifyRemoteCertificate,omitempty" yaml:"verifyRemoteCertificate,omitempty"`
	AuthHeaderFrom          *SecretRef `json:"authHeaderFrom,omitempty" yaml:"authHeaderFrom,omitempty"`
}

// ReplicationPolicy describes a Harbor replication policy using registry names instead of IDs.
type ReplicationPolicy struct {
	Name                    string              `json:"name" yaml:"name"`
	Description             *string             `json:"description,omitempty" yaml:"description,omitempty"`
	Enabled                 *bool               `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Mode                    string              `json:"mode,omitempty" yaml:"mode,omitempty"`
	Registry                string              `json:"registry,omitempty" yaml:"registry,omitempty"`
	DestinationNamespace    *string             `json:"destinationNamespace,omitempty" yaml:"destinationNamespace,omitempty"`
	DestinationReplaceCount *int8               `json:"destinationNamespaceReplaceCount,omitempty" yaml:"destinationNamespaceReplaceCount,omitempty"`
	Override                *bool               `json:"override,omitempty" yaml:"override,omitempty"`
	ReplicateDeletion       *bool               `json:"replicateDeletion,omitempty" yaml:"replicateDeletion,omitempty"`
	CopyByChunk             *bool               `json:"copyByChunk,omitempty" yaml:"copyByChunk,omitempty"`
	Speed                   *int32              `json:"speed,omitempty" yaml:"speed,omitempty"`
	Filters                 []ReplicationFilter `json:"filters,omitempty" yaml:"filters,omitempty"`
	Trigger                 *ReplicationTrigger `json:"trigger,omitempty" yaml:"trigger,omitempty"`
}

// ReplicationFilter limits resources selected by a replication policy.
type ReplicationFilter struct {
	Type       string `json:"type" yaml:"type"`
	Value      string `json:"value,omitempty" yaml:"value,omitempty"`
	Decoration string `json:"decoration,omitempty" yaml:"decoration,omitempty"`
}

// ReplicationTrigger describes when a replication policy runs.
type ReplicationTrigger struct {
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	Cron string `json:"cron,omitempty" yaml:"cron,omitempty"`
}

// NewConfiguration returns an empty document using the current schema version.
func NewConfiguration() *Configuration {
	return &Configuration{APIVersion: APIVersion, Kind: Kind}
}

// Validate checks schema identity, resource identity, and local references.
func (c *Configuration) Validate() error {
	return c.validate(true)
}

func (c *Configuration) validate(requireComplete bool) error {
	if c == nil {
		return fmt.Errorf("configuration is empty")
	}
	if c.APIVersion != APIVersion {
		return fmt.Errorf("unsupported apiVersion %q; expected %q", c.APIVersion, APIVersion)
	}
	if c.Kind != Kind {
		return fmt.Errorf("unsupported kind %q; expected %q", c.Kind, Kind)
	}
	if err := validateNames("registry", c.Spec.Registries, func(value Registry) string { return value.Name }); err != nil {
		return err
	}
	if err := validateNames("project", c.Spec.Projects, func(value Project) string { return value.Name }); err != nil {
		return err
	}
	if err := validateNames("replication policy", c.Spec.ReplicationPolicies, func(value ReplicationPolicy) string { return value.Name }); err != nil {
		return err
	}

	for _, registry := range c.Spec.Registries {
		if registry.Credential == nil {
			continue
		}
		for _, ref := range []*SecretRef{registry.Credential.AccessKeyFrom, registry.Credential.AccessSecretFrom} {
			if ref != nil && strings.TrimSpace(ref.Env) == "" {
				return fmt.Errorf("registry %q has an empty credential environment reference", registry.Name)
			}
		}
	}
	for _, project := range c.Spec.Projects {
		if err := validateNames("webhook in project "+project.Name, project.Webhooks, func(value Webhook) string { return value.Name }); err != nil {
			return err
		}
		for _, webhook := range project.Webhooks {
			for i, target := range webhook.Targets {
				if requireComplete && (target.NotifyType == nil || strings.TrimSpace(*target.NotifyType) == "") {
					return fmt.Errorf("webhook %q in project %q target at index %d has no notifyType", webhook.Name, project.Name, i)
				}
				if requireComplete && (target.Endpoint == nil || strings.TrimSpace(*target.Endpoint) == "") {
					return fmt.Errorf("webhook %q in project %q target at index %d has no endpoint", webhook.Name, project.Name, i)
				}
				if target.AuthHeaderFrom != nil && strings.TrimSpace(target.AuthHeaderFrom.Env) == "" {
					return fmt.Errorf("webhook %q in project %q target at index %d has an empty auth-header environment reference", webhook.Name, project.Name, i)
				}
			}
		}
	}
	for _, policy := range c.Spec.ReplicationPolicies {
		if policy.Mode != "" && policy.Mode != "push" && policy.Mode != "pull" {
			return fmt.Errorf("replication policy %q mode must be push or pull", policy.Name)
		}
		if requireComplete && policy.Mode == "" {
			return fmt.Errorf("replication policy %q must define a mode", policy.Name)
		}
		if requireComplete && strings.TrimSpace(policy.Registry) == "" {
			return fmt.Errorf("replication policy %q must reference a registry", policy.Name)
		}
		for i, filter := range policy.Filters {
			if strings.TrimSpace(filter.Type) == "" {
				return fmt.Errorf("replication policy %q filter at index %d has no type", policy.Name, i)
			}
		}
		if requireComplete && policy.Trigger != nil && strings.TrimSpace(policy.Trigger.Type) == "" {
			return fmt.Errorf("replication policy %q trigger has no type", policy.Name)
		}
	}
	return nil
}

func validateNames[T any](kind string, values []T, name func(T) string) error {
	seen := make(map[string]struct{}, len(values))
	for i, item := range values {
		value := strings.TrimSpace(name(item))
		if value == "" {
			return fmt.Errorf("%s at index %d has no name", kind, i)
		}
		if _, ok := seen[value]; ok {
			return fmt.Errorf("duplicate %s %q", kind, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

// Normalize sorts resources so repeated exports produce stable output.
func (c *Configuration) Normalize() {
	if c == nil {
		return
	}
	slices.SortFunc(c.Spec.Registries, func(a, b Registry) int { return cmp.Compare(a.Name, b.Name) })
	slices.SortFunc(c.Spec.Projects, func(a, b Project) int { return cmp.Compare(a.Name, b.Name) })
	for i := range c.Spec.Projects {
		project := &c.Spec.Projects[i]
		slices.SortFunc(project.Webhooks, func(a, b Webhook) int { return cmp.Compare(a.Name, b.Name) })
		for j := range project.Webhooks {
			slices.Sort(project.Webhooks[j].Events)
		}
	}
	slices.SortFunc(c.Spec.ReplicationPolicies, func(a, b ReplicationPolicy) int {
		return cmp.Compare(a.Name, b.Name)
	})
	for i := range c.Spec.ReplicationPolicies {
		policy := &c.Spec.ReplicationPolicies[i]
		slices.SortFunc(policy.Filters, func(a, b ReplicationFilter) int {
			return cmp.Or(
				cmp.Compare(a.Type, b.Type),
				cmp.Compare(a.Value, b.Value),
				cmp.Compare(a.Decoration, b.Decoration),
			)
		})
	}
}
