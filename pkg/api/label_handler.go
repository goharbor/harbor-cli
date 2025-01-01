package api

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/label"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/label/create"
)

func CreateLabel(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	_, err = client.Label.CreateLabel(ctx, &label.CreateLabelParams{Label: &models.Label{Name: opts.Name, Color: opts.Color, Description: opts.Description, Scope: opts.Scope}})

	if err != nil {
		return err
	}

	fmt.Printf("Label %s created\n", opts.Name)
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

	fmt.Println("label deleted successfully")

	return nil
}

func ListLabel(opts ...ListFlags) (*label.ListLabelsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags

	if len(opts) > 0 {
		listFlags = opts[0]
	}
	scope := "g"
	response, err := client.Label.ListLabels(ctx, &label.ListLabelsParams{
		Page:      &listFlags.Page,
		PageSize:  &listFlags.PageSize,
		Q:         &listFlags.Q,
		Sort:      &listFlags.Sort,
		Scope:     &scope,
		ProjectID: &listFlags.ProjectID,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func UpdateLabel(updateView *models.Label, Labelid int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	labelUpdate := &models.Label{
		Name:        updateView.Name,
		Color:       updateView.Color,
		Description: updateView.Description,
		Scope:       updateView.Scope,
	}

	_, err = client.Label.UpdateLabel(
		ctx,
		&label.UpdateLabelParams{LabelID: Labelid, Label: labelUpdate},
	)
	if err != nil {
		return err
	}

	fmt.Println("label updated successfully")

	return nil
}

func GetLabel(labelid int64) *models.Label {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil
	}
	response, err := client.Label.GetLabelByID(ctx, &label.GetLabelByIDParams{LabelID: labelid})
	if err != nil {
		return nil
	}

	return response.GetPayload()
}

func GetLabelIdByName(labelName string) (int64, error) {
	var opts ListFlags

	l, err := ListLabel(opts)
	if err != nil {
		return 0, err
	}

	for _, label := range l.Payload {
		if label.Name == labelName {
			return label.ID, nil
		}
	}

	return 0, err
}
