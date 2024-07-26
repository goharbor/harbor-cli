package registry

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/registry/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewListRegistryCommand creates a new `harbor list registry` command
func ListRegistryCommand() *cobra.Command {
	var opts api.ListFlags
	var formatFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list registry",
		Run: func(cmd *cobra.Command, args []string) {
			response, err := api.ListRegistries(opts)
			if err != nil {
				log.Fatalf("failed to get projects list: %v", err)
			}

			if formatFlag != "" {
				if formatFlag == "json" {
					utils.PrintPayloadInJSONFormat(response)
				} else if formatFlag == "yaml" {
					utils.PrintPayloadInYAMLFormat(response)
				} else {
					log.Errorf("invalid output format: %s", formatFlag)
				}
			} else {
				list.ListRegistries(response.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")
	flags.StringVarP(&formatFlag, "output-format", "o", "", "Output format. One of: json|yaml")

	return cmd
}
