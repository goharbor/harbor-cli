package labels

import (
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/label/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UpdateLableCommand() *cobra.Command {
	var opts models.Label

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update labels",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var labelId int64
			updateflags := &api.ListFlags{
				Scope: opts.Scope,
			}

			if len(args) > 0 {
				labelId, err = strconv.ParseInt(args[0], 10, 64)
			} else {
				labelId = prompt.GetLabelIdFromUser(*updateflags)
			}
			if err != nil {
				log.Errorf("failed to parse label id: %v", err)
			}

			opts = *api.GetLabel(labelId)
			updateView := &models.Label{
				Name:        opts.Name,
				Color:       opts.Color,
				Description: opts.Description,
				Scope:       opts.Scope,
			}

			update.UpdateLabelView(updateView)
			err = api.UpdateLabel(updateView, labelId)
			if err != nil {
				log.Errorf("failed to update label: %v", err)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "default(global).p for project labels.Query scope of the label")

	return cmd
}
