package systeminfo

import (
	"github.com/spf13/cobra"

)

func SystemInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "systeminfo",
		Short: "Interact with system information",
		Long:  `Commands to interact with the Harbor system information, including general info, volumes, and certificates.`,
	}

	cmd.AddCommand(
		GetInfoCommand(),
		GetVolumesCommand(),
		GetCertCommand(),
	)

	return cmd
}