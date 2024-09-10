package cveallowlist

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/systemcve/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list system level allowlist of cve",
		Run: func(cmd *cobra.Command, args []string) {
			cve, err := api.ListSystemCve()
			if err != nil {
				log.Fatalf("failed to get system cve list: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(cve)
				return
			}

			list.ListSystemCve(cve.Payload)
		},
	}

	return cmd
}
