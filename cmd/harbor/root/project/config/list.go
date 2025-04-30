package config

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ListConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list [NAME|ID]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name")
			} else {
				projectNameOrID := args[0]

				response, err := api.ListConfig(isID, projectNameOrID)
				if err != nil {
					log.Errorf("failed to view metadata: %v", err)
				} else {
					utils.PrintPayloadInJSONFormat(response.Payload)
				}
			}

		},
	}

	return cmd
}
