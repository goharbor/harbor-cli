package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/label"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/goharbor/harbor-cli/pkg/views/label/create"
	lview "github.com/goharbor/harbor-cli/pkg/views/label/select"
)

func GetLabelIdFromUser(opts *models.Label) int64 {
	labelId := make(chan int64)
	go func() {
		ctx, client, _ := utils.ContextWithClient()
		response, err := client.Label.ListLabels(ctx, &label.ListLabelsParams{Scope: &opts.Scope, ProjectID: &opts.ProjectID})
		if err != nil {
			log.Fatal(err)
		}
		lview.LabelList(response.Payload, labelId)

	}()

	return <-labelId
}

func CreateLabels(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Label.CreateLabel(ctx, &label.CreateLabelParams{Label: &models.Label{Name: opts.Name, Color: opts.Color, Description: opts.Description, Scope: opts.Scope}})

	if err != nil {
		return err
	}

	log.Infof("Label %s created", opts.Name)
	return nil
}

func DeleteLabel(Labelid int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Label.DeleteLabel(ctx, &label.DeleteLabelParams{LabelID: Labelid})

	if err != nil {
		return err
	}

	log.Info("label deleted successfully")

	return nil
}

func ListLabel(opts ...ListLabelFlags) (*label.ListLabelsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListLabelFlags

	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Label.ListLabels(ctx, &label.ListLabelsParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Q:        &listFlags.Q,
		Sort:     &listFlags.Sort,
		Scope: &listFlags.Scope,
		ProjectID: &listFlags.ProjectID,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
