package project

import (
	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// GetProjectCommand creates a new `harbor get project` command
func GetProjectCommand() *cobra.Command {
	var opts config.GetProjectOptions

	cmd := &cobra.Command{
		Use:   "project [NAME|ID]",
		Short: "get project by name or id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.ProjectNameOrID = args[0]
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return api.RunGetProject(opts, credentialName, config.OutputType, config.WideOutput)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&config.OutputType, "output", "o", "", "Output type [json/yaml]")
	return cmd
}