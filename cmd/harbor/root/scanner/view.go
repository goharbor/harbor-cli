package scanner

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ViewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "get scanner by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registrationID := args[0]
				err = api.GetScanner(registrationID)
			} else {
				registrationID := prompt.GetScannerIdFromUser()
				err = api.GetScanner(registrationID)
			}

			if err != nil {
				log.Errorf("failed to get scanner: %v", err)
			}

		},
	}

	return cmd
}
