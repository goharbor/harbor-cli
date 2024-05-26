package tags

import "github.com/spf13/cobra"

func ArtifactTagsCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "tags",
		Short:   "Manage tags of an artifact",
		Example: ` harbor artifact tags list <project>/<repository>/<reference>`,
	}

	cmd.AddCommand(
		ListTagsCommand(),
	)

	return cmd
}
