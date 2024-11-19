package artifact

import (
	"github.com/spf13/cobra"
)

func Artifact() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "artifact",
		Short:   "Manage artifacts",
		Long:    `Manage artifacts in Harbor Repository`,
		Example: `  harbor artifact list`,
	}

	cmd.AddCommand(
		ListArtifactCommand(),
		ViewArtifactCommmand(),
		DeleteArtifactCommand(),
		ScanArtifactCommand(),
		ArtifactTagsCmd(),
	)

	return cmd
}
