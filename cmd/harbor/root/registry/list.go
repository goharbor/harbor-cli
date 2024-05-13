package registry

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type listRegistryOptions struct {
	page     int64
	pageSize int64
	q        string
	sort     string
}

// NewListRegistryCommand creates a new `harbor list registry` command
func ListRegistryCommand() *cobra.Command {
	var opts listRegistryOptions

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return runListRegistry(opts, credentialName)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.pageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.sort, "sort", "", "", "Sort the resource list in ascending or descending order")

	return cmd
}

func runListRegistry(opts listRegistryOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Registry.ListRegistries(ctx, &registry.ListRegistriesParams{Page: &opts.page, PageSize: &opts.pageSize, Q: &opts.q, Sort: &opts.sort})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.GetPayload())
	return nil
}
