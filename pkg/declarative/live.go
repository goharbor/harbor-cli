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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/configure"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/quota"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/replication"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/webhook"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

const exportPageSize int64 = 100

// LiveBackend reconciles configuration through Harbor's v2 API.
type LiveBackend struct {
	client *v2client.HarborAPI

	registries map[string]*models.Registry
	projects   map[string]*models.Project
	quotas     map[string]*models.Quota
	webhooks   map[string]map[string]*models.WebhookPolicy
	policies   map[string]*models.ReplicationPolicy
}

// NewLiveBackend creates a backend using the active Harbor CLI context.
func NewLiveBackend() (*LiveBackend, error) {
	client, err := utils.GetClient()
	if err != nil {
		return nil, err
	}
	return &LiveBackend{client: client}, nil
}

// Snapshot exports all supported API-managed resources and refreshes apply-time indexes.
func (b *LiveBackend) Snapshot(ctx context.Context) (*Configuration, error) {
	b.registries = make(map[string]*models.Registry)
	b.projects = make(map[string]*models.Project)
	b.quotas = make(map[string]*models.Quota)
	b.webhooks = make(map[string]map[string]*models.WebhookPolicy)
	b.policies = make(map[string]*models.ReplicationPolicy)

	document := NewConfiguration()
	configuration, err := b.client.Configure.GetConfigurations(ctx, &configure.GetConfigurationsParams{})
	if err != nil {
		return nil, fmt.Errorf("get system configuration: %w", err)
	}
	document.Spec.System = systemConfiguration(configuration.Payload)

	registries, err := b.listRegistries(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range registries {
		b.registries[item.Name] = item
		document.Spec.Registries = append(document.Spec.Registries, registryFromModel(item))
	}

	quotas, err := b.listQuotas(ctx)
	if err != nil {
		return nil, err
	}
	quotaByProjectID := make(map[int64]*models.Quota, len(quotas))
	for _, item := range quotas {
		projectID, ok := quotaProjectID(item.Ref)
		if ok {
			quotaByProjectID[projectID] = item
		}
	}

	projects, err := b.listProjects(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range projects {
		b.projects[item.Name] = item
		projectDocument := projectFromModel(item, b.registries)
		if projectQuota := quotaByProjectID[int64(item.ProjectID)]; projectQuota != nil {
			projectDocument.Quota = maps.Clone(projectQuota.Hard)
			b.quotas[item.Name] = projectQuota
		}

		webhooks, listErr := b.client.Webhook.ListWebhookPoliciesOfProject(ctx, &webhook.ListWebhookPoliciesOfProjectParams{
			ProjectNameOrID: item.Name,
		})
		if listErr != nil {
			return nil, fmt.Errorf("list webhooks for project %q: %w", item.Name, listErr)
		}
		b.webhooks[item.Name] = make(map[string]*models.WebhookPolicy, len(webhooks.Payload))
		for _, policy := range webhooks.Payload {
			b.webhooks[item.Name][policy.Name] = policy
			projectDocument.Webhooks = append(projectDocument.Webhooks, webhookFromModel(policy))
		}
		document.Spec.Projects = append(document.Spec.Projects, projectDocument)
	}

	policies, err := b.listReplicationPolicies(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range policies {
		policy, conversionErr := replicationPolicyFromModel(item)
		if conversionErr != nil {
			return nil, conversionErr
		}
		b.policies[item.Name] = item
		document.Spec.ReplicationPolicies = append(document.Spec.ReplicationPolicies, policy)
	}

	document.Normalize()
	return document, nil
}

// Apply executes one planned action.
func (b *LiveBackend) Apply(ctx context.Context, desired *Configuration, action Action) error {
	switch action.Resource {
	case resourceRegistry:
		value, ok := findByName(desired.Spec.Registries, action.Name, func(item Registry) string { return item.Name })
		if !ok {
			return fmt.Errorf("registry %q is absent from desired state", action.Name)
		}
		return b.applyRegistry(ctx, value, action.Operation)
	case resourceProject:
		value, ok := findByName(desired.Spec.Projects, action.Name, func(item Project) string { return item.Name })
		if !ok {
			return fmt.Errorf("project %q is absent from desired state", action.Name)
		}
		return b.applyProject(ctx, value, action.Operation)
	case resourceQuota:
		value, ok := findByName(desired.Spec.Projects, action.Name, func(item Project) string { return item.Name })
		if !ok {
			return fmt.Errorf("project %q is absent from desired state", action.Name)
		}
		return b.applyQuota(ctx, value)
	case resourceWebhook:
		projectValue, ok := findByName(desired.Spec.Projects, action.Project, func(item Project) string { return item.Name })
		if !ok {
			return fmt.Errorf("project %q is absent from desired state", action.Project)
		}
		value, ok := findByName(projectValue.Webhooks, action.Name, func(item Webhook) string { return item.Name })
		if !ok {
			return fmt.Errorf("webhook %q is absent from project %q", action.Name, action.Project)
		}
		return b.applyWebhook(ctx, action.Project, value, action.Operation)
	case resourceReplicationPolicy:
		value, ok := findByName(desired.Spec.ReplicationPolicies, action.Name, func(item ReplicationPolicy) string { return item.Name })
		if !ok {
			return fmt.Errorf("replication policy %q is absent from desired state", action.Name)
		}
		return b.applyReplicationPolicy(ctx, value, action.Operation)
	case resourceSystem:
		return b.applySystem(ctx, desired.Spec.System)
	default:
		return fmt.Errorf("unsupported resource type %q", action.Resource)
	}
}

func (b *LiveBackend) listRegistries(ctx context.Context) ([]*models.Registry, error) {
	var result []*models.Registry
	for page := int64(1); ; page++ {
		response, err := b.client.Registry.ListRegistries(ctx, &registry.ListRegistriesParams{Page: &page, PageSize: ptr(exportPageSize)})
		if err != nil {
			return nil, fmt.Errorf("list registries: %w", err)
		}
		result = append(result, response.Payload...)
		if len(response.Payload) < int(exportPageSize) {
			return result, nil
		}
	}
}

func (b *LiveBackend) listProjects(ctx context.Context) ([]*models.Project, error) {
	var result []*models.Project
	for page := int64(1); ; page++ {
		response, err := b.client.Project.ListProjects(ctx, &project.ListProjectsParams{Page: &page, PageSize: ptr(exportPageSize)})
		if err != nil {
			return nil, fmt.Errorf("list projects: %w", err)
		}
		result = append(result, response.Payload...)
		if len(response.Payload) < int(exportPageSize) {
			return result, nil
		}
	}
}

func (b *LiveBackend) listQuotas(ctx context.Context) ([]*models.Quota, error) {
	var result []*models.Quota
	for page := int64(1); ; page++ {
		response, err := b.client.Quota.ListQuotas(ctx, &quota.ListQuotasParams{Page: &page, PageSize: ptr(exportPageSize)})
		if err != nil {
			return nil, fmt.Errorf("list quotas: %w", err)
		}
		result = append(result, response.Payload...)
		if len(response.Payload) < int(exportPageSize) {
			return result, nil
		}
	}
}

func (b *LiveBackend) listReplicationPolicies(ctx context.Context) ([]*models.ReplicationPolicy, error) {
	var result []*models.ReplicationPolicy
	for page := int64(1); ; page++ {
		response, err := b.client.Replication.ListReplicationPolicies(ctx, &replication.ListReplicationPoliciesParams{Page: &page, PageSize: ptr(exportPageSize)})
		if err != nil {
			return nil, fmt.Errorf("list replication policies: %w", err)
		}
		result = append(result, response.Payload...)
		if len(response.Payload) < int(exportPageSize) {
			return result, nil
		}
	}
}

func (b *LiveBackend) applyRegistry(ctx context.Context, desired Registry, operation Operation) error {
	var credential *models.RegistryCredential
	if desired.Credential != nil {
		var err error
		credential, err = resolveRegistryCredential(desired.Credential)
		if err != nil {
			return err
		}
	}
	if operation == OperationCreate {
		if desired.Type == nil || desired.URL == nil {
			return fmt.Errorf("creating registry %q requires type and url", desired.Name)
		}
		request := &models.Registry{Name: desired.Name, Type: *desired.Type, URL: *desired.URL, Credential: credential}
		if desired.Description != nil {
			request.Description = *desired.Description
		}
		if desired.Insecure != nil {
			request.Insecure = *desired.Insecure
		}
		if _, err := b.client.Registry.CreateRegistry(ctx, &registry.CreateRegistryParams{Registry: request}); err != nil {
			return err
		}
		return b.refreshRegistry(ctx, desired.Name)
	}

	current := b.registries[desired.Name]
	if current == nil {
		return fmt.Errorf("registry %q not found", desired.Name)
	}
	if desired.Type != nil && *desired.Type != current.Type {
		return fmt.Errorf("registry %q type is immutable (current %q, desired %q)", desired.Name, current.Type, *desired.Type)
	}
	request := &models.RegistryUpdate{
		Name:        &desired.Name,
		Description: desired.Description,
		URL:         desired.URL,
		Insecure:    desired.Insecure,
	}
	if desired.Credential != nil {
		if desired.Credential.Type != nil {
			request.CredentialType = &credential.Type
		}
		if desired.Credential.AccessKeyFrom != nil {
			request.AccessKey = &credential.AccessKey
		}
		if desired.Credential.AccessSecretFrom != nil {
			request.AccessSecret = &credential.AccessSecret
		}
	}
	_, err := b.client.Registry.UpdateRegistry(ctx, &registry.UpdateRegistryParams{ID: current.ID, Registry: request})
	return err
}

func (b *LiveBackend) refreshRegistry(ctx context.Context, name string) error {
	page := int64(1)
	response, err := b.client.Registry.ListRegistries(ctx, &registry.ListRegistriesParams{Name: &name, Page: &page, PageSize: ptr(exportPageSize)})
	if err != nil {
		return fmt.Errorf("refresh registry %q: %w", name, err)
	}
	for _, item := range response.Payload {
		if item.Name == name {
			b.registries[name] = item
			return nil
		}
	}
	return fmt.Errorf("registry %q was created but could not be read back", name)
}

func (b *LiveBackend) applyProject(ctx context.Context, desired Project, operation Operation) error {
	request, err := b.projectRequest(desired)
	if err != nil {
		return err
	}
	if operation == OperationCreate {
		request.ProjectName = desired.Name
		if desired.Quota != nil {
			if storage, ok := desired.Quota["storage"]; ok {
				request.StorageLimit = &storage
			}
		}
		if _, err := b.client.Project.CreateProject(ctx, &project.CreateProjectParams{Project: request}); err != nil {
			return err
		}
		return b.refreshProject(ctx, desired.Name)
	}

	resourceName := true
	_, err = b.client.Project.UpdateProject(ctx, &project.UpdateProjectParams{
		ProjectNameOrID: desired.Name,
		XIsResourceName: &resourceName,
		Project:         request,
	})
	return err
}

func (b *LiveBackend) projectRequest(desired Project) (*models.ProjectReq, error) {
	request := &models.ProjectReq{Public: desired.Public}
	current := b.projects[desired.Name]
	request.Metadata = metadataToModel(desired.Metadata, metadataOf(current))
	if desired.Public != nil {
		if request.Metadata == nil {
			request.Metadata = metadataToModel(&ProjectMetadata{}, metadataOf(current))
		}
		request.Metadata.Public = strconv.FormatBool(*desired.Public)
	}
	if desired.Registry != nil {
		registryModel := b.registries[*desired.Registry]
		if registryModel == nil {
			return nil, fmt.Errorf("project %q references unknown registry %q", desired.Name, *desired.Registry)
		}
		request.RegistryID = &registryModel.ID
	}
	return request, nil
}

func (b *LiveBackend) refreshProject(ctx context.Context, name string) error {
	resourceName := true
	response, err := b.client.Project.GetProject(ctx, &project.GetProjectParams{ProjectNameOrID: name, XIsResourceName: &resourceName})
	if err != nil {
		return fmt.Errorf("refresh project %q: %w", name, err)
	}
	b.projects[name] = response.Payload
	return nil
}

func (b *LiveBackend) applyQuota(ctx context.Context, desired Project) error {
	quotaModel := b.quotas[desired.Name]
	if quotaModel == nil {
		if err := b.refreshQuota(ctx, desired.Name); err != nil {
			return err
		}
		quotaModel = b.quotas[desired.Name]
	}
	_, err := b.client.Quota.UpdateQuota(ctx, &quota.UpdateQuotaParams{
		ID:   quotaModel.ID,
		Hard: &models.QuotaUpdateReq{Hard: maps.Clone(desired.Quota)},
	})
	return err
}

func (b *LiveBackend) refreshQuota(ctx context.Context, projectName string) error {
	projectModel := b.projects[projectName]
	if projectModel == nil {
		if err := b.refreshProject(ctx, projectName); err != nil {
			return err
		}
		projectModel = b.projects[projectName]
	}
	quotas, err := b.listQuotas(ctx)
	if err != nil {
		return err
	}
	for _, item := range quotas {
		projectID, ok := quotaProjectID(item.Ref)
		if ok && projectID == int64(projectModel.ProjectID) {
			b.quotas[projectName] = item
			return nil
		}
	}
	return fmt.Errorf("quota for project %q not found", projectName)
}

func (b *LiveBackend) applyWebhook(ctx context.Context, projectName string, desired Webhook, operation Operation) error {
	current := b.webhooks[projectName][desired.Name]
	policy, err := webhookToModel(desired, current)
	if err != nil {
		return err
	}
	if operation == OperationCreate {
		_, err = b.client.Webhook.CreateWebhookPolicyOfProject(ctx, &webhook.CreateWebhookPolicyOfProjectParams{
			ProjectNameOrID: projectName,
			Policy:          policy,
		})
		return err
	}
	if current == nil {
		return fmt.Errorf("webhook %q not found in project %q", desired.Name, projectName)
	}
	_, err = b.client.Webhook.UpdateWebhookPolicyOfProject(ctx, &webhook.UpdateWebhookPolicyOfProjectParams{
		ProjectNameOrID: projectName,
		WebhookPolicyID: current.ID,
		Policy:          policy,
	})
	return err
}

func (b *LiveBackend) applyReplicationPolicy(ctx context.Context, desired ReplicationPolicy, operation Operation) error {
	current := b.policies[desired.Name]
	policy, err := b.replicationPolicyToModel(desired, current)
	if err != nil {
		return err
	}
	if operation == OperationCreate {
		_, err = b.client.Replication.CreateReplicationPolicy(ctx, &replication.CreateReplicationPolicyParams{Policy: policy})
		return err
	}
	if current == nil {
		return fmt.Errorf("replication policy %q not found", desired.Name)
	}
	_, err = b.client.Replication.UpdateReplicationPolicy(ctx, &replication.UpdateReplicationPolicyParams{ID: current.ID, Policy: policy})
	return err
}

func (b *LiveBackend) replicationPolicyToModel(desired ReplicationPolicy, current *models.ReplicationPolicy) (*models.ReplicationPolicy, error) {
	policy := &models.ReplicationPolicy{Name: desired.Name}
	if current != nil {
		*policy = *current
	}
	policy.Name = desired.Name
	registryModel := b.registries[desired.Registry]
	if registryModel == nil {
		return nil, fmt.Errorf("replication policy %q references unknown registry %q", desired.Name, desired.Registry)
	}
	if desired.Mode == "pull" {
		policy.SrcRegistry = registryModel
		policy.DestRegistry = nil
	} else {
		policy.SrcRegistry = nil
		policy.DestRegistry = registryModel
	}
	assign(desired.Description, &policy.Description)
	assign(desired.Enabled, &policy.Enabled)
	assign(desired.DestinationNamespace, &policy.DestNamespace)
	if desired.DestinationReplaceCount != nil {
		policy.DestNamespaceReplaceCount = desired.DestinationReplaceCount
	}
	assign(desired.Override, &policy.Override)
	if desired.ReplicateDeletion != nil {
		policy.ReplicateDeletion = *desired.ReplicateDeletion
		policy.Deletion = *desired.ReplicateDeletion
	}
	if desired.CopyByChunk != nil {
		policy.CopyByChunk = desired.CopyByChunk
	}
	if desired.Speed != nil {
		policy.Speed = desired.Speed
	}
	if desired.Filters != nil {
		policy.Filters = make([]*models.ReplicationFilter, 0, len(desired.Filters))
		for _, filter := range desired.Filters {
			policy.Filters = append(policy.Filters, &models.ReplicationFilter{Type: filter.Type, Value: filter.Value, Decoration: filter.Decoration})
		}
	}
	if desired.Trigger != nil {
		policy.Trigger = &models.ReplicationTrigger{Type: desired.Trigger.Type}
		if desired.Trigger.Cron != "" {
			policy.Trigger.TriggerSettings = &models.ReplicationTriggerSettings{Cron: desired.Trigger.Cron}
		}
	}
	return policy, nil
}

func (b *LiveBackend) applySystem(ctx context.Context, desired map[string]any) error {
	encoded, err := json.Marshal(desired)
	if err != nil {
		return fmt.Errorf("encode system configuration: %w", err)
	}
	var configuration models.Configurations
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&configuration); err != nil {
		return fmt.Errorf("decode system configuration: %w", err)
	}
	_, err = b.client.Configure.UpdateConfigurations(ctx, &configure.UpdateConfigurationsParams{Configurations: &configuration})
	return err
}

func systemConfiguration(response *models.ConfigurationsResponse) map[string]any {
	result := make(map[string]any)
	if response == nil {
		return result
	}
	value := reflect.ValueOf(response).Elem()
	typeInfo := value.Type()
	for i := range value.NumField() {
		field := value.Field(i)
		if field.Kind() != reflect.Pointer || field.IsNil() {
			continue
		}
		jsonName := strings.Split(typeInfo.Field(i).Tag.Get("json"), ",")[0]
		if jsonName == "" || jsonName == "-" || isSecretSystemField(jsonName) {
			continue
		}
		item := field.Elem()
		valueField := item.FieldByName("Value")
		if valueField.IsValid() && valueField.CanInterface() {
			result[jsonName] = valueField.Interface()
		}
	}
	return result
}

func isSecretSystemField(name string) bool {
	switch name {
	case "ldap_search_password", "oidc_client_secret", "uaa_client_secret":
		return true
	default:
		return false
	}
}

func registryFromModel(value *models.Registry) Registry {
	return Registry{
		Name:        value.Name,
		Type:        ptr(value.Type),
		URL:         ptr(value.URL),
		Description: ptr(value.Description),
		Insecure:    ptr(value.Insecure),
	}
}

func projectFromModel(value *models.Project, registries map[string]*models.Registry) Project {
	result := Project{Name: value.Name, Metadata: metadataFromModel(value.Metadata)}
	if value.Metadata != nil {
		result.Public = parseBool(value.Metadata.Public)
	}
	if result.Public == nil {
		result.Public = ptr(false)
	}
	if value.RegistryID != 0 {
		for name, registryModel := range registries {
			if registryModel.ID == value.RegistryID {
				result.Registry = &name
				break
			}
		}
	}
	return result
}

func metadataFromModel(value *models.ProjectMetadata) *ProjectMetadata {
	if value == nil {
		return nil
	}
	return &ProjectMetadata{
		AutoScan:                 parseBoolPointer(value.AutoScan),
		AutoSBOMGeneration:       parseBoolPointer(value.AutoSbomGeneration),
		EnableContentTrust:       parseBoolPointer(value.EnableContentTrust),
		EnableContentTrustCosign: parseBoolPointer(value.EnableContentTrustCosign),
		PreventVulnerableImages:  parseBoolPointer(value.PreventVul),
		ProxySpeedKB:             parseIntPointer(value.ProxySpeedKb),
		ReuseSystemCVEAllowlist:  parseBoolPointer(value.ReuseSysCVEAllowlist),
		Severity:                 clonePointer(value.Severity),
	}
}

func metadataToModel(desired *ProjectMetadata, current *models.ProjectMetadata) *models.ProjectMetadata {
	if desired == nil {
		return nil
	}
	result := &models.ProjectMetadata{}
	if current != nil {
		*result = *current
	}
	assignBoolString(desired.AutoScan, &result.AutoScan)
	assignBoolString(desired.AutoSBOMGeneration, &result.AutoSbomGeneration)
	assignBoolString(desired.EnableContentTrust, &result.EnableContentTrust)
	assignBoolString(desired.EnableContentTrustCosign, &result.EnableContentTrustCosign)
	assignBoolString(desired.PreventVulnerableImages, &result.PreventVul)
	assignIntString(desired.ProxySpeedKB, &result.ProxySpeedKb)
	assignBoolString(desired.ReuseSystemCVEAllowlist, &result.ReuseSysCVEAllowlist)
	if desired.Severity != nil {
		result.Severity = desired.Severity
	}
	return result
}

func metadataOf(projectModel *models.Project) *models.ProjectMetadata {
	if projectModel == nil {
		return nil
	}
	return projectModel.Metadata
}

func webhookFromModel(value *models.WebhookPolicy) Webhook {
	result := Webhook{
		Name:        value.Name,
		Description: ptr(value.Description),
		Enabled:     ptr(value.Enabled),
		Events:      slices.Clone(value.EventTypes),
	}
	for _, target := range value.Targets {
		if target == nil {
			continue
		}
		result.Targets = append(result.Targets, WebhookTarget{
			NotifyType:              ptr(target.Type),
			Endpoint:                ptr(target.Address),
			PayloadFormat:           ptr(string(target.PayloadFormat)),
			VerifyRemoteCertificate: ptr(!target.SkipCertVerify),
		})
	}
	return result
}

func webhookToModel(desired Webhook, current *models.WebhookPolicy) (*models.WebhookPolicy, error) {
	policy := &models.WebhookPolicy{Name: desired.Name, Enabled: true}
	if current != nil {
		*policy = *current
	}
	policy.Name = desired.Name
	assign(desired.Description, &policy.Description)
	assign(desired.Enabled, &policy.Enabled)
	if desired.Events != nil {
		policy.EventTypes = slices.Clone(desired.Events)
	}
	if desired.Targets == nil {
		if current == nil {
			return nil, fmt.Errorf("creating webhook %q requires at least one target", desired.Name)
		}
		return policy, nil
	}
	if len(desired.Targets) == 0 {
		return nil, fmt.Errorf("webhook %q requires at least one target", desired.Name)
	}
	policy.Targets = make([]*models.WebhookTargetObject, 0, len(desired.Targets))
	for i, desiredTarget := range desired.Targets {
		target := &models.WebhookTargetObject{}
		assign(desiredTarget.NotifyType, &target.Type)
		assign(desiredTarget.Endpoint, &target.Address)
		if desiredTarget.PayloadFormat != nil {
			target.PayloadFormat = models.PayloadFormatType(*desiredTarget.PayloadFormat)
		}
		if desiredTarget.VerifyRemoteCertificate != nil {
			target.SkipCertVerify = !*desiredTarget.VerifyRemoteCertificate
		}
		if desiredTarget.AuthHeaderFrom != nil {
			value, err := resolveSecret(desiredTarget.AuthHeaderFrom)
			if err != nil {
				return nil, fmt.Errorf("resolve webhook target %d auth header: %w", i, err)
			}
			target.AuthHeader = value
		}
		if target.Type == "" || target.Address == "" {
			return nil, fmt.Errorf("webhook %q target %d requires notifyType and endpoint", desired.Name, i)
		}
		policy.Targets = append(policy.Targets, target)
	}
	return policy, nil
}

func replicationPolicyFromModel(value *models.ReplicationPolicy) (ReplicationPolicy, error) {
	result := ReplicationPolicy{
		Name:                    value.Name,
		Description:             ptr(value.Description),
		Enabled:                 ptr(value.Enabled),
		DestinationNamespace:    ptr(value.DestNamespace),
		DestinationReplaceCount: clonePointer(value.DestNamespaceReplaceCount),
		Override:                ptr(value.Override),
		ReplicateDeletion:       ptr(value.ReplicateDeletion),
		CopyByChunk:             clonePointer(value.CopyByChunk),
		Speed:                   clonePointer(value.Speed),
	}
	registryModel := value.DestRegistry
	result.Mode = "push"
	if value.SrcRegistry != nil {
		result.Mode = "pull"
		registryModel = value.SrcRegistry
	}
	if registryModel == nil || registryModel.Name == "" {
		return ReplicationPolicy{}, fmt.Errorf("replication policy %q has no external registry", value.Name)
	}
	result.Registry = registryModel.Name
	for _, filter := range value.Filters {
		result.Filters = append(result.Filters, ReplicationFilter{
			Type:       filter.Type,
			Value:      fmt.Sprint(filter.Value),
			Decoration: filter.Decoration,
		})
	}
	if value.Trigger != nil {
		result.Trigger = &ReplicationTrigger{Type: value.Trigger.Type}
		if value.Trigger.TriggerSettings != nil {
			result.Trigger.Cron = value.Trigger.TriggerSettings.Cron
		}
	}
	return result, nil
}

func resolveRegistryCredential(value *RegistryCredential) (*models.RegistryCredential, error) {
	result := &models.RegistryCredential{}
	if value.Type != nil {
		result.Type = *value.Type
	}
	var err error
	if value.AccessKeyFrom != nil {
		result.AccessKey, err = resolveSecret(value.AccessKeyFrom)
		if err != nil {
			return nil, fmt.Errorf("resolve registry access key: %w", err)
		}
	}
	if value.AccessSecretFrom != nil {
		result.AccessSecret, err = resolveSecret(value.AccessSecretFrom)
		if err != nil {
			return nil, fmt.Errorf("resolve registry access secret: %w", err)
		}
	}
	return result, nil
}

func resolveSecret(ref *SecretRef) (string, error) {
	value, ok := os.LookupEnv(ref.Env)
	if !ok {
		return "", fmt.Errorf("environment variable %q is not set", ref.Env)
	}
	return value, nil
}

func quotaProjectID(ref any) (int64, bool) {
	encoded, err := json.Marshal(ref)
	if err != nil {
		return 0, false
	}
	var value struct {
		ID int64 `json:"id"`
	}
	if json.Unmarshal(encoded, &value) != nil || value.ID == 0 {
		return 0, false
	}
	return value.ID, true
}

func parseBoolPointer(value *string) *bool {
	if value == nil {
		return nil
	}
	return parseBool(*value)
}

func parseBool(value string) *bool {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseIntPointer(value *string) *int64 {
	if value == nil {
		return nil
	}
	parsed, err := strconv.ParseInt(*value, 10, 64)
	if err != nil {
		return nil
	}
	return &parsed
}

func assignBoolString(source *bool, destination **string) {
	if source != nil {
		value := strconv.FormatBool(*source)
		*destination = &value
	}
}

func assignIntString(source *int64, destination **string) {
	if source != nil {
		value := strconv.FormatInt(*source, 10)
		*destination = &value
	}
}

func assign[T any](source *T, destination *T) {
	if source != nil {
		*destination = *source
	}
}

func clonePointer[T any](value *T) *T {
	if value == nil {
		return nil
	}
	return ptr(*value)
}

func findByName[T any](values []T, name string, key func(T) string) (T, bool) {
	for _, value := range values {
		if key(value) == name {
			return value, true
		}
	}
	var zero T
	return zero, false
}

func ptr[T any](value T) *T {
	return &value
}
