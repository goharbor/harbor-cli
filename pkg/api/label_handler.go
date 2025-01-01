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
		return fmt.Errorf("failed to initialize client context for label creation")
	}

	_, err = client.Label.CreateLabel(ctx, &label.CreateLabelParams{
		Label: &models.Label{
			Name:        opts.Name,
			Color:       opts.Color,
			Description: opts.Description,
			Scope:       opts.Scope,
		},
	})

	if err != nil {
		switch err.(type) {
		case *label.CreateLabelBadRequest:
			return fmt.Errorf("invalid request while creating label: %s", opts.Name)
		case *label.CreateLabelConflict:
			return fmt.Errorf("label already exists: %s", opts.Name)
		case *label.CreateLabelUnsupportedMediaType:
			return fmt.Errorf("unsupported media type for label creation")
		case *label.CreateLabelInternalServerError:
			return fmt.Errorf("internal server error occurred during label creation")
		case *label.CreateLabelUnauthorized:
			return fmt.Errorf("unauthorized access to create label")
		default:
			return fmt.Errorf("unknown error occurred while creating label: %v", err)
		}
	}

	fmt.Printf("Label '%s' created successfully\n", opts.Name)
	return nil
}

func DeleteLabel(Labelid int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client context for label deletion")
	}
	_, err = client.Label.DeleteLabel(ctx, &label.DeleteLabelParams{LabelID: Labelid})

	if err != nil {
		switch err.(type) {
		case *label.DeleteLabelNotFound:
			return fmt.Errorf("label with ID %d not found", Labelid)
		case *label.DeleteLabelBadRequest:
			return fmt.Errorf("invalid request while deleting label with ID %d", Labelid)
		case *label.DeleteLabelInternalServerError:
			return fmt.Errorf("internal server error occurred during label deletion")
		case *label.DeleteLabelUnauthorized:
			return fmt.Errorf("unauthorized access to delete label with ID %d", Labelid)
		default:
			return fmt.Errorf("unknown error occurred while deleting label with ID %d: %v", Labelid, err)
		}
	}

	fmt.Println("Label deleted successfully")

	return nil
}

func ListLabel(opts ...ListFlags) (*label.ListLabelsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client context for listing labels")
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
		switch err.(type) {
		case *label.ListLabelsBadRequest:
			return nil, fmt.Errorf("invalid request while listing labels for the project")
		case *label.ListLabelsInternalServerError:
			return nil, fmt.Errorf("internal server error occurred while listing labels")
		case *label.ListLabelsUnauthorized:
			return nil, fmt.Errorf("unauthorized access to list labels")
		default:
			return nil, fmt.Errorf("unknown error occurred while listing labels: %v", err)
		}
	}

	return response, nil
}

func UpdateLabel(updateView *models.Label, Labelid int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return fmt.Errorf("failed to initialize client context for label update")
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
		switch err.(type) {
		case *label.UpdateLabelNotFound:
			return fmt.Errorf("label with ID %d not found", Labelid)
		case *label.UpdateLabelBadRequest:
			return fmt.Errorf("invalid request while updating label with ID %d", Labelid)
		case *label.UpdateLabelInternalServerError:
			return fmt.Errorf("internal server error occurred during label update")
		case *label.UpdateLabelUnauthorized:
			return fmt.Errorf("unauthorized access to update label with ID %d", Labelid)
		default:
			return fmt.Errorf("unknown error occurred while updating label with ID %d: %v", Labelid, err)
		}
	}

	fmt.Println("Label updated successfully")

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
		return 0, fmt.Errorf("failed to list labels: %v", err)
	}

	for _, label := range l.Payload {
		if label.Name == labelName {
			return label.ID, nil
		}
	}

	return 0, fmt.Errorf("label with name '%s' not found", labelName)
}
