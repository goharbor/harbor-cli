package root

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/registry"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "harbor",
		Short:        "Official Harbor CLI",
		SilenceUsage: true,
		Long:         "Official Harbor CLI",
		Example: `
// Base command:
harbor

// Display help about the command:
harbor help
`,
		// RunE: func(cmd *cobra.Command, args []string) error {

		// },
	}

	root.AddCommand(
		versionCommand(),
		LoginCommand(),
		project.Project(),
		registry.Registry(),
	)

	return root
}
