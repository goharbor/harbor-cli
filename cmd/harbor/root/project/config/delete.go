package config

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var isID bool

func DeleteConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete [NAME|ID] ...[KEY]",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name or id and the meta names to delete")
			} else if len(args) == 1 {
				fmt.Println("Please provide the meta names")
			} else {
				projectNameOrID := args[0]
				metaNames := make([]string, 0)
				for i := 1; i < len(args); i++ {
					metaNames = append(metaNames, args[i])
				}

				err := api.DeleteConfig(isID, projectNameOrID, metaNames)
				if err != nil {
					log.Errorf("failed to delete metadata: %v", err)
				}
			}

		},
	}

	return cmd
}
