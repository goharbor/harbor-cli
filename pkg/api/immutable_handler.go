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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/immutable"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/immutable/create"
	log "github.com/sirupsen/logrus"
)

func CreateImmutable(opts create.CreateView, projectName string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	tagSelector := &models.ImmutableSelector{
		Decoration: opts.TagSelectors.Decoration,
		Pattern:    opts.TagSelectors.Pattern,
	}
	scope := models.ImmutableSelector{
		Decoration: opts.ScopeSelectors.Decoration,
		Pattern:    opts.ScopeSelectors.Pattern,
	}
	scopeSelector := map[string][]models.ImmutableSelector{
		"repository": {
			scope,
		},
	}

	_, err = client.Immutable.CreateImmuRule(ctx, &immutable.CreateImmuRuleParams{ProjectNameOrID: projectName, ImmutableRule: &models.ImmutableRule{TagSelectors: []*models.ImmutableSelector{tagSelector}, ScopeSelectors: scopeSelector}})

	if err != nil {
		return err
	}

	log.Info("Added Tag Immutability Rule")
	return nil
}

func ListImmutable(projectName string) (immutable.ListImmuRulesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return immutable.ListImmuRulesOK{}, err
	}
	response, err := client.Immutable.ListImmuRules(ctx, &immutable.ListImmuRulesParams{ProjectNameOrID: projectName})
	if err != nil {
		return immutable.ListImmuRulesOK{}, err
	}

	return *response, nil
}

func DeleteImmutable(projectName string, ImmutableID int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Immutable.DeleteImmuRule(ctx, &immutable.DeleteImmuRuleParams{ProjectNameOrID: projectName, ImmutableRuleID: ImmutableID})
	if err != nil {
		return err
	}

	log.Info("immutable rule deleted successfully")

	return nil
}
