package labels

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteLabelCommand() *cobra.Command {
	var opts models.Label
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete label",
		Example: "harbor label delete [labelname]",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			deleteView := &api.ListFlags{
				Scope: opts.Scope,
			}

			if len(args) > 0 {
				labelId, _ := api.GetLabelIdByName(args[0])
				err = api.DeleteLabel(labelId)
			} else {
				labelId := prompt.GetLabelIdFromUser(*deleteView)
				err = api.DeleteLabel(labelId)
			}
			if err != nil {
				log.Errorf("failed to delete label: %v", err)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "default(global).'p' for project labels.Query scope of the label")

	return cmd
}
