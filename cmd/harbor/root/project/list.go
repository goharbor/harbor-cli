package project

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/project/list"
	plist "github.com/goharbor/harbor-cli/pkg/views/project/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewListProjectCommand creates a new `harbor list project` command
func ListProjectCommand() *cobra.Command {
	var opts plist.ListProjectOptions

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list project",
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := api.ListProject(opts)
			if err != nil {
				log.Fatalf("failed to get projects list: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(projects)
				return
			}

			list.ListProjects(projects.Payload)
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
	flags.BoolVarP(&opts.WithDetail, "with-detail", "", true, "Bool value indicating whether return detailed information of the project")

	return cmd
}
