package scan_all

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/spf13/cobra"
)

func StopScanAllCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop scanning all artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return api.StopScanAll()
		},
	}

	return cmd
}
