package labels

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/label/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListLabelCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list labels",
		Run: func(cmd *cobra.Command, args []string) {
			label, err := api.ListLabel(opts)
			if err != nil {
				log.Fatalf("failed to get label list: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(label)
				return
			}
			list.ListLabels(label.Payload)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 20, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "default(global).'p' for project labels.Query scope of the label")
	flags.Int64VarP(&opts.ProjectID, "projectid", "i", 1, "project ID when query project labels")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the label list in ascending or descending order")

	return cmd
}
