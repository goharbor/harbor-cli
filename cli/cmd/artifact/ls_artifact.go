package artifact

import (
	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// NewListArtifactCommand creates a new `harbor list artifact` command
func ListArtifactCommand() *cobra.Command {
	var opts config.ListArtifactOptions

	cmd := &cobra.Command{
		Use:   "artifact",
		Short: "list artifact",
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, _ := cmd.Flags().GetString(constants.CredentialNameOption)

			return api.RunListArtifact(opts, credentialName, config.OutputType, config.WideOutput)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.ProjectName, "projectname", "", "", "Name of the project")
	flags.StringVarP(&opts.RepositoryName, "reponame", "", "", "Name of the repository")
	flags.StringVarP(&config.OutputType, "output", "o", "", "Output type [json/yaml]")
	flags.BoolVarP(&config.WideOutput, "owide", "", false, "Wide output result [true/false]")
	return cmd
}