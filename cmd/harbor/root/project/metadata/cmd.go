package metadata

import "github.com/spf13/cobra"

var (
	isID bool
)

func Metadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Manage project metadata",
	}
	cmd.AddCommand(
		AddMetadataCommand(),
		DeleteMetadataCommand(),
		ViewMetadataCommand(),
		UpdateMetadataCommand(),
		ListMetadataCommand(),
	)

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "Use project ID instead of name")

	return cmd
}
