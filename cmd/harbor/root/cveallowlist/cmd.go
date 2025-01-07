package cveallowlist

import (
	"github.com/spf13/cobra"
)

func CVEAllowlist() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cve-allowlist",
		Short:   "Manage system CVE allowlist",
		Long:    `Managing CVE lists that are intentionally excluded from vulnerability scanning`,
		Example: `harbor cve-allowlist list`,
	}
	cmd.AddCommand(
		ListCveCommand(),
		AddCveAllowlistCommand(),
	)

	return cmd
}
