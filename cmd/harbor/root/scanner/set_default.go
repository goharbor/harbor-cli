package scanner

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func SetDefaultCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-default",
		Short:   "set default scanner",
		Aliases: []string{"sd"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error

			if len(args) > 0 {
				registrationID := args[0]
				err = api.SetDefaultScanner(registrationID)
			} else {
				registrationID := prompt.GetScannerIdFromUser()
				err = api.SetDefaultScanner(registrationID)
			}

			if err != nil {
				log.Errorf("failed to set default scanner: %v", err)
			}

		},
	}

	return cmd
}
