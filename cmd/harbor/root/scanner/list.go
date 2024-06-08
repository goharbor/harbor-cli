package scanner

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func ListScannerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List scanners",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			scanners, err := api.ListScanners()
			if err != nil {
				cmd.PrintErrf("failed to list scanners: %v", err)
				return
			}

			utils.PrintPayloadInJSONFormat(scanners.Payload)
		},
	}

	return cmd
}
