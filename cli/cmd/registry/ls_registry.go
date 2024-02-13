package registry

import (
	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// ListRegistryCommand creates a new `harbor list registry` command
func ListRegistryCommand() *cobra.Command {
	var opts config.ListRegistryOptions

	cmd := &cobra.Command{
		Use:   "registry",
		Short: "list registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return api.RunListRegistry(opts, credentialName, config.OutputType, config.WideOutput)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")
	flags.StringVarP(&config.OutputType, "output", "o", "", "Output type [json/yaml]")
	flags.BoolVarP(&config.WideOutput, "wide", "", false, "Wide output result [true/false]")
	return cmd
}
