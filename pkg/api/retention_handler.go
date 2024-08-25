package api

import (
	"strconv"

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
        Decoration:  	opts.TagSelectors.Decoration,
        Pattern: 		opts.TagSelectors.Pattern,
		Extras: 		opts.TagSelectors.Extras,
    }
	scope := models.RetentionSelector{
        Decoration: 	opts.ScopeSelectors.Decoration,
        Pattern: 		opts.ScopeSelectors.Pattern,
    }
	scopeSelector := map[string][]models.RetentionSelector{
        "repository": {
			scope,
		},
    }
	param := make(map[string] interface{})
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
		Action: 		opts.Action,
		ScopeSelectors: scopeSelector,
    	TagSelectors: 	[]*models.RetentionSelector{tagSelector},
    	Template: 		opts.Template,
    	Params:   		param,
	})

	triggerSettings := map[string]string{
		"cron": "",
	}

	_, err = client.Retention.CreateRetention(ctx, &retention.CreateRetentionParams{Policy: &models.RetentionPolicy{Scope: &models.RetentionPolicyScope{Level: opts.Scope.Level,Ref: int64(projectId)},Trigger: &models.RetentionRuleTrigger{Kind: models.ScheduleObjTypeSchedule,Settings: triggerSettings},Algorithm: opts.Algorithm,Rules: rule}})
	if err != nil {
		return err
	}

	log.Info("Added Tag Retention Rule")
	return nil
}