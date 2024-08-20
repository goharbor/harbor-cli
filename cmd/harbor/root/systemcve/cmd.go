package systemcve

import (
	"github.com/spf13/cobra"
)

func SystemCVEAllowlist() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "systemcve",
		Short:   "Manage system CVE allowlist",
		Long:    `Manage system level allowlist of CVE`,
		Example: `harbor systemcve list`,
	}
	cmd.AddCommand(
		ListCveCommand(),
		UpdateSystemCveCommand(),
	)

	return cmd
}
