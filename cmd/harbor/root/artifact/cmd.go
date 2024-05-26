package artifact

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/artifact/tags"
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
		InfoArtifactCommmand(),
		DeleteArtifactCommand(),
		ScanArtifactCommand(),
		tags.ArtifactTagsCmd(),
	)

	return cmd
}
