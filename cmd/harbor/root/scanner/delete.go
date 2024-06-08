package scanner

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete scanner",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registrationID := args[0]
				err = api.DeleteScanner(registrationID)
			} else {
				registrationID := prompt.GetScannerIdFromUser()
				err = api.DeleteScanner(registrationID)
			}

			if err != nil {
				log.Errorf("failed to delete scanner: %v", err)
			} else {
				log.Infof("scanner deleted successfully")
			}

		},
	}

	return cmd
}
