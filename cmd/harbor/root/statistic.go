package root

import (
	"log"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/stat/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// versionCommand represents the version command
func StatisticCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stats",
		Short:   "statistics of projects and repositories",
		Long:    `Get the statistic information about the projects and repositories`,
		Example: `  harbor stats`,
		Run: func(cmd *cobra.Command, args []string) {
			stats, err := api.GetStats()
			if err != nil {
				log.Fatal(err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag == "json" {
				utils.PrintPayloadInJSONFormat(stats)
				return
			} else if FormatFlag == "wide" {
				list.ListStatistics(stats.Payload, true)
			} else {
				list.ListStatistics(stats.Payload, false)
			}
		},
	}

	return cmd
}
