package root

import (
	"log"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/info/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Lists the info of the Harbor system
func InfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Short:   "Get general system info",
		Example: `  harbor info`,
		Run: func(cmd *cobra.Command, args []string) {
			generalInfo, err := api.GetSystemInfo()
			if err != nil {
				log.Fatal(err)
			}

			stats, err := api.GetStats()
			if err != nil {
				log.Fatal(err)
			}

			sysVolume, err := api.GetSystemVolumes()
			if err != nil {
				log.Fatal(err)
			}

			// CreateSystemInfo
			systemInfo := list.CreateSystemInfo(generalInfo.Payload, stats.Payload, sysVolume.Payload)

			FormatFlag := viper.GetString("output-format")
			if FormatFlag == "json" {
				utils.PrintPayloadInJSONFormat(systemInfo)
				return
			}

			list.ListInfo(&systemInfo)
		},
	}

	return cmd
}
