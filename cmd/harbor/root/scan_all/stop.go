package scan_all

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func StopScanAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop scanning all artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			logrus.Info("Stopping scan all operation")
			err := api.StopScanAll()
			if err != nil {
				logrus.Errorf("Failed to stop scan all operation: %v", utils.ParseHarborErrorMsg(err))
				return err
			}
			logrus.Info("Successfully stopped scan all operation")
			return nil
		},
	}

	return cmd
}
