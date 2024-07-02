package schedule

import (
	"github.com/spf13/cobra"
)

func Schedule() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Schedule jobs in Harbor",
	}
	cmd.AddCommand(
		ListScheduleCommand(),
	)

	return cmd
}
