package scan_all

import "github.com/spf13/cobra"

func ScanAll() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scan-all",
		Short: "Scan all artifacts",
	}

	cmd.AddCommand(
		UpdateScanAllScheduleCommand(),
		StopScanAllCommand(),
		ViewScanAllScheduleCommand(),
		GetScanAllMetricsCommand(),
		RunScanAllCommand(),
	)

	return cmd
}
