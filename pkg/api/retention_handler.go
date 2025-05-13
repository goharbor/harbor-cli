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
package api

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/retention"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/retention/create"
	log "github.com/sirupsen/logrus"
)

func CreateRetention(opts create.CreateView, projectIDorName string, isName bool) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	newRule, err := generateRetentionRule(opts)
	if err != nil {
		return err
	}
	retentionIDStr := ""
	retentionIDStr, err = GetRetentionId(projectIDorName, isName)
	if err != nil && err.Error() != "No retention policy exists for this project" {
		return err
	}

	if retentionIDStr != "" {
		return UpdateRetention(retentionIDStr, newRule)
	}

	triggerSettings := map[string]string{
		"cron": "",
	}
	var projectID int
	if isName {
		project, _ := client.Project.GetProject(ctx, &project.GetProjectParams{
			XIsResourceName: &isName,
			ProjectNameOrID: projectIDorName,
		})
		projectID = int(project.Payload.ProjectID)
	} else {
		projectID, err = strconv.Atoi(projectIDorName)
		if err != nil {
			return fmt.Errorf("failed to convert project ID to int: %w", err)
		}
	}
	_, err = client.Retention.CreateRetention(ctx, &retention.CreateRetentionParams{
		Policy: &models.RetentionPolicy{
			Scope: &models.RetentionPolicyScope{
				Level: opts.Scope.Level,
				Ref:   int64(projectID),
			},
			Trigger: &models.RetentionRuleTrigger{
				Kind:     models.ScheduleObjTypeSchedule,
				Settings: triggerSettings,
			},
			Algorithm: opts.Algorithm,
			Rules:     []*models.RetentionRule{newRule},
		},
	})

	if err != nil {
		return err
	}

	log.Info("Created new Tag Retention Policy")
	return nil
}

func ListRetention(retentionID string) (retention.GetRetentionOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return retention.GetRetentionOK{}, err
	}
	retentionIDint, err := strconv.Atoi(retentionID)
	if err != nil {
		return retention.GetRetentionOK{}, err
	}
	response, err := client.Retention.GetRetention(ctx, &retention.GetRetentionParams{ID: int64(retentionIDint)})
	if err != nil {
		return retention.GetRetentionOK{}, err
	}
	return *response, nil
}

func GetRetentionId(projectNameorID string, isName bool) (string, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return "", err
	}

	response, err := client.Project.GetProject(ctx, &project.GetProjectParams{
		XIsResourceName: &isName,
		ProjectNameOrID: projectNameorID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get project with such name or ID: %w", err)
	}

	if response.Payload.Metadata == nil || response.Payload.Metadata.RetentionID == nil {
		return "", errors.New("No retention policy exists for this project")
	}
	retentionid := *response.Payload.Metadata.RetentionID

	return retentionid, nil
}

func DeleteRetention(projectName string, ruleIndex int) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	retentionIDStr, err := GetRetentionId(projectName, true)
	if err != nil {
		return err
	}

	retentionResp, err := ListRetention(retentionIDStr)
	if err != nil {
		return err
	}

	existingPolicy := retentionResp.Payload
	if ruleIndex < 0 || ruleIndex >= len(existingPolicy.Rules) {
		return fmt.Errorf("invalid rule index")
	}

	existingPolicy.Rules = append(existingPolicy.Rules[:ruleIndex], existingPolicy.Rules[ruleIndex+1:]...)

	_, err = client.Retention.UpdateRetention(ctx, &retention.UpdateRetentionParams{
		ID:     int64(retentionResp.Payload.ID),
		Policy: existingPolicy,
	})
	if err != nil {
		return err
	}

	log.Infof("Deleted rule at index %d from retention policy", ruleIndex)
	return nil
}

func UpdateRetention(retentionIDStr string, newRule *models.RetentionRule) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	retentionResp, err := ListRetention(retentionIDStr)
	if err != nil {
		return err
	}

	existingPolicy := retentionResp.Payload
	if len(existingPolicy.Rules) >= 15 {
		return fmt.Errorf("cannot add rule: retention policy already has 15 rules")
	}
	existingPolicy.Rules = append(existingPolicy.Rules, newRule)

	_, err = client.Retention.UpdateRetention(ctx, &retention.UpdateRetentionParams{
		ID:     int64(retentionResp.Payload.ID),
		Policy: existingPolicy,
	})
	if err != nil {
		return err
	}

	log.Info("Updated existing retention policy")
	return nil
}

func generateRetentionRule(opts create.CreateView) (*models.RetentionRule, error) {
	tagSelector := &models.RetentionSelector{
		Decoration: opts.TagSelectors.Decoration,
		Pattern:    opts.TagSelectors.Pattern,
		Extras:     opts.TagSelectors.Extras,
	}
	scope := models.RetentionSelector{
		Decoration: opts.ScopeSelectors.Decoration,
		Pattern:    opts.ScopeSelectors.Pattern,
	}
	scopeSelector := map[string][]models.RetentionSelector{
		"repository": {scope},
	}
	param := make(map[string]interface{})
	if opts.Template != "always" {
		value, err := strconv.Atoi(opts.Params.Value)
		if err != nil {
			return nil, err
		}
		param[opts.Template] = value
	} else {
		param = nil
	}

	return &models.RetentionRule{
		Action:         opts.Action,
		ScopeSelectors: scopeSelector,
		TagSelectors:   []*models.RetentionSelector{tagSelector},
		Template:       opts.Template,
		Params:         param,
	}, nil
}
