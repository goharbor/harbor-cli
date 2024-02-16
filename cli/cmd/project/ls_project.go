package project

import (
	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// NewListProjectCommand creates a new `harbor list project` command
func ListProjectCommand() *cobra.Command {
	var opts config.ListProjectOptions

	cmd := &cobra.Command{
		Use:   "project",
		Short: "list project",
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return api.RunListProject(opts, credentialName, config.OutputType, config.WideOutput)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the project")
	flags.StringVarP(&opts.Owner, "owner", "", "", "Name of the project owner")
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.BoolVarP(&opts.Public, "public", "", true, "Project is public or private")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")
	flags.StringVarP(&config.OutputType, "output", "o", "", "Output type [json/yaml]")
	flags.BoolVarP(&config.WideOutput, "wide", "", false, "Wide output result [true/false]")
	return cmd
}