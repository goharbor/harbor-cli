package labels

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/label/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateLableCommand() *cobra.Command {
	opts := &models.Label{}

	cmd := &cobra.Command{
		Use:     "update",
		Short:   "update label",
		Example: "harbor label update [labelname]",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var labelId int64
			updateflags := api.ListFlags{
				Scope: opts.Scope,
			}

			if len(args) > 0 {
				labelId, err = api.GetLabelIdByName(args[0])
			} else {
				labelId = prompt.GetLabelIdFromUser(updateflags)
			}
			if err != nil {
				log.Errorf("failed to parse label id: %v", err)
			}

			existingLabel := api.GetLabel(labelId)
			if existingLabel == nil {
				log.Errorf("label is not found")
				return
			}
			updateView := &models.Label{
				Name:        existingLabel.Name,
				Color:       existingLabel.Color,
				Description: existingLabel.Description,
				Scope:       existingLabel.Scope,
			}

			flags := cmd.Flags()
			if flags.Changed("name") {
				updateView.Name = opts.Name
			}
			if flags.Changed("color") {
				updateView.Color = opts.Color
			}
			if flags.Changed("description") {
				updateView.Description = opts.Description
			}
			if flags.Changed("scope") {
				updateView.Scope = opts.Scope
			}

			update.UpdateLabelView(updateView)
			err = api.UpdateLabel(updateView, labelId)
			if err != nil {
				log.Errorf("failed to update label: %v", err)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "Name of the label")
	flags.StringVarP(&opts.Color, "color", "", "", "Color of the label.color is in hex value")
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "Scope of the label. eg- g(global), p(specific project)")
	flags.StringVarP(&opts.Description, "description", "d", "", "Description of the label")

	return cmd
}
