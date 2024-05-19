package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/project/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ListProjectOptions struct {
	Name       string
	Owner      string
	Page       int64
	PageSize   int64
	Public     bool
	Q          string
	Sort       string
	WithDetail bool
}

// NewListProjectCommand creates a new `harbor list project` command
func ListProjectCommand() *cobra.Command {
	var opts ListProjectOptions

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list project",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			credentialName := viper.GetString("current-credential-name")
			client := utils.GetClientByCredentialName(credentialName)
			projects, err := RunListProject(opts, ctx, client.Project)
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

func RunListProject(opts ListProjectOptions, ctx context.Context, projectInterface ProjectInterface) (*project.ListProjectsOK, error) {

	response, err := projectInterface.ListProjects(ctx, &project.ListProjectsParams{Name: &opts.Name, Owner: &opts.Owner, Page: &opts.Page, PageSize: &opts.PageSize, Public: &opts.Public, Q: &opts.Q, Sort: &opts.Sort, WithDetail: &opts.WithDetail})
	if err != nil {
		return nil, err
	}
	return response, nil
}
