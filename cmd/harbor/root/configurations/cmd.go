package configurations

import (
	"github.com/spf13/cobra"
)

func ConfigurationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Manage system configurations",
		Long:    "Manage system configurations",
		Example: `harbor config get`,
	}

	cmd.AddCommand(
		GetConfigCmd(),
		UpdateConfigCmd(),
	)

	return cmd
}
