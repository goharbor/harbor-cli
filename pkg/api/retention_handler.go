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

func CreateRetention(opts create.CreateView, projectId int32) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

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
		"repository": {
			scope,
		},
	}
	param := make(map[string]interface{})
	if opts.Template == "always" {
		param = nil
	} else {
		value, err := strconv.Atoi(opts.Params.Value)
		if err != nil {
			return err
		}
		param[opts.Template] = value
	}

	var rule []*models.RetentionRule
	rule = append(rule, &models.RetentionRule{
		Action:         opts.Action,
		ScopeSelectors: scopeSelector,
		TagSelectors:   []*models.RetentionSelector{tagSelector},
		Template:       opts.Template,
		Params:         param,
	})

	triggerSettings := map[string]string{
		"cron": "",
	}

	_, err = client.Retention.CreateRetention(ctx, &retention.CreateRetentionParams{Policy: &models.RetentionPolicy{Scope: &models.RetentionPolicyScope{Level: opts.Scope.Level, Ref: int64(projectId)}, Trigger: &models.RetentionRuleTrigger{Kind: models.ScheduleObjTypeSchedule, Settings: triggerSettings}, Algorithm: opts.Algorithm, Rules: rule}})
	if err != nil {
		return err
	}

	log.Info("Added Tag Retention Rule")
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
		return "", fmt.Errorf("failed to get project: %w", err)
	}

	if response.Payload.Metadata == nil || response.Payload.Metadata.RetentionID == nil {
		return "", errors.New("No retention policy exists for this project")
	}
	retentionid := *response.Payload.Metadata.RetentionID

	return retentionid, nil
}

func DeleteRetention(retentionID string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	retentionIDint, err := strconv.Atoi(retentionID)
	if err != nil {
		return err
	}
	_, err = client.Retention.DeleteRetention(ctx, &retention.DeleteRetentionParams{ID: int64(retentionIDint)})
	if err != nil {
		return err
	}

	log.Info("retention rule deleted successfully")

	return nil
}
