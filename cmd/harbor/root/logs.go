package root

import (
	"log"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/logs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Logs() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get recent logs of the projects which the user is a member of",
		Run: func(cmd *cobra.Command, args []string) {
			logs, err := api.AuditLogs(opts)
			if err != nil {
				log.Fatalf("failed to get projects list: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(logs.Payload)
				return
			}
			list.ListLogs(logs.Payload)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(
		&opts.Sort,
		"sort",
		"",
		"",
		"Sort the resource list in ascending or descending order",
	)

	return cmd
}
