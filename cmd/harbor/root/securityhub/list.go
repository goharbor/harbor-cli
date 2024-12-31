package securityhub

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/securityhub/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListVulnerabilityCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show all the vulnerability list",
		Run: func(cmd *cobra.Command, args []string) {
			vulnerability, err := api.ListVulnerability(opts)
			if err != nil {
				log.Fatalf("failed to get vulnerability list: %v", err)
				return
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(vulnerability, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.ListVulnerability(vulnerability.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")

	return cmd
}
