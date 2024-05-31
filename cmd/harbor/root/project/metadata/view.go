package metadata

import (
	"fmt"
	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ViewMetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "view [NAME|ID] [KEY]",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name or id and the meta name")
			} else if len(args) == 1 {
				fmt.Println("Please provide the meta name")
			} else {
				projectNameOrID := args[0]
				metaName := args[1]

				err := api.ViewMetadata(isID, projectNameOrID, metaName)
				if err != nil {
					log.Errorf("failed to view metadata: %v", err)
				}
			}

		},
	}

	return cmd
}
