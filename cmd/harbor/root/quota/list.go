package quota

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/quota/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Lists the Quotas specified for each project
func ListQuotaCommand() *cobra.Command {
	var opts api.ListQuotaFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list quotas",
		Run: func(cmd *cobra.Command, args []string) {
			quota, err := api.ListQuota(opts)
			if err != nil {
				log.Fatalf("failed to get projects list: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(quota)
				return
			}

			list.ListQuotas(quota.Payload)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.Reference, "ref", "", "", "Reference type of quota")
	flags.StringVarP(&opts.ReferenceID, "refid", "", "", "Reference ID of quota")
	flags.StringVarP(
		&opts.Sort,
		"sort",
		"",
		"",
		"Sort the resource list in ascending or descending order",
	)

	return cmd
}
