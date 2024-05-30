package metadata

import "github.com/spf13/cobra"

func Metadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Manage project metadata",
	}
	cmd.AddCommand(
		AddMetadataCommand(),
		DeleteMetadataCommand(),
		ViewMetadataCommand(),
	)

	return cmd
}
