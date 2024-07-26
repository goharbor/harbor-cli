package project

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/project/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewListProjectCommand creates a new `harbor list project` command
func ListProjectCommand() *cobra.Command {
	var opts api.ListFlags
	var formatFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list project",
		Run: func(cmd *cobra.Command, args []string) {
			response, err := api.ListProject(opts)
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
				list.ListProjects(response.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "", "", "Name of the project")
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.BoolVarP(&opts.Public, "public", "", false, "Project is public or private")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")
	flags.StringVarP(&formatFlag, "output-format", "o", "", "Output format. One of: json|yaml")

	return cmd
}
