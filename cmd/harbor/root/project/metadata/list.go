package metadata

import (
	"fmt"
	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ListMetadataCommand() *cobra.Command {
	var isID bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list [NAME|ID]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name")
			} else {
				projectNameOrID := args[0]

				err := api.ListMetadata(isID, projectNameOrID)
				if err != nil {
					log.Errorf("failed to view metadata: %v", err)
				}
			}

		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "Use project ID instead of name")

	return cmd
}
