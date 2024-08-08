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
        Decoration:  opts.TagSelectors.Decoration,
        Pattern: opts.TagSelectors.Pattern,
    }
	scope := models.ImmutableSelector{
        Decoration: opts.ScopeSelectors.Decoration,
        Pattern: opts.ScopeSelectors.Pattern,
    }
	scopeSelector := map[string][]models.ImmutableSelector{
        "repository": {
			scope,
		},
    }

	_, err = client.Immutable.CreateImmuRule(ctx, &immutable.CreateImmuRuleParams{ProjectNameOrID: projectName,ImmutableRule: &models.ImmutableRule{TagSelectors: []*models.ImmutableSelector{tagSelector},ScopeSelectors: scopeSelector}})

	if err != nil {
		return err
	}

	log.Info("Added Tag Immutability Rule")
	return nil
}

func ListImmutable(projectName string)(immutable.ListImmuRulesOK,error){
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return immutable.ListImmuRulesOK{}, err
	}
	response,err := client.Immutable.ListImmuRules(ctx,&immutable.ListImmuRulesParams{ProjectNameOrID: projectName})
	if err != nil {
		return immutable.ListImmuRulesOK{}, err
	}

	return *response, nil
}

func DeleteImmutable(projectName string,ImmutableID int64) error{
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Immutable.DeleteImmuRule(ctx,&immutable.DeleteImmuRuleParams{ProjectNameOrID: projectName,ImmutableRuleID: ImmutableID})
	if err != nil {
		return err
	}

	log.Info("immutable rule deleted successfully")

	return nil
}