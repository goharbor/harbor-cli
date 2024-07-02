package labels

import (
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteLabelCommand() *cobra.Command {
	var opts models.Label
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete label by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			deleteView := &models.Label{
				Scope: opts.Scope,
				ProjectID: opts.ProjectID,
			}

			if len(args) > 0 {
				labelId, _ := strconv.ParseInt(args[0], 10, 64)
				err = api.DeleteLabel(labelId)
			} else {
				labelId := api.GetLabelIdFromUser(deleteView)
				err = api.DeleteLabel(labelId)
			}
			if err != nil {
				log.Errorf("failed to delete label: %v", err)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "default(global).p for project labels.Query scope of the label")
	flags.Int64VarP(&opts.ProjectID, "projectid", "i", 1, "project ID when query project labels")

	return cmd
}