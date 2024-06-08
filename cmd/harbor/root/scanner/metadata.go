package scanner

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func MetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "get scanner metadata by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registrationID := args[0]
				err = api.GetScannerMetadata(registrationID)
			} else {
				registrationID := prompt.GetScannerIdFromUser()
				err = api.GetScannerMetadata(registrationID)
			}

			if err != nil {
				log.Errorf("failed to get scanner metadata: %v", err)
			}

		},
	}

	return cmd
}
