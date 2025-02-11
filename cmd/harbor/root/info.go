package root

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/info/list"
	log "github.com/sirupsen/logrus"
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
			if FormatFlag != "" {
				err = utils.PrintFormat(systemInfo, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.ListInfo(&systemInfo)
			}
		},
	}

	return cmd
}
