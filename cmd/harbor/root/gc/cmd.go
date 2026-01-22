package gc

import (
	"github.com/spf13/cobra"
)

func GC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gc",
		Short: "Manage Garbage Collection",
		Long:  "Manage Garbage Collection in Harbor (schedule, history, logs)",
	}

	cmd.AddCommand(
		ListGCCommand(),
		ViewGCScheduleCommand(),
		UpdateGCScheduleCommand(),
		GetGCLogCommand(),
		RunGCCommand(),
	)

	return cmd
}
