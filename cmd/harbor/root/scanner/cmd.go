package scanner

import "github.com/spf13/cobra"

func Scanner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scanner",
		Short: "scanner commands",
	}

	cmd.AddCommand(
		CreateScannerCommand(),
		ListScannerCommand(),
		ViewCommand(),
		MetadataCommand(),
		SetDefaultCommand(),
		DeleteCommand(),
		UpdateCommand(),
	)

	return cmd
}
