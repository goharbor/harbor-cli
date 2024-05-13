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

type listProjectOptions struct {
	name       string
	owner      string
	page       int64
	pageSize   int64
	public     bool
	q          string
	sort       string
	withDetail bool
}

// NewListProjectCommand creates a new `harbor list project` command
func ListProjectCommand() *cobra.Command {
	var opts listProjectOptions

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list project",
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := RunListProject(opts)
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
	flags.StringVarP(&opts.name, "name", "", "", "Name of the project")
	flags.StringVarP(&opts.owner, "owner", "", "", "Name of the project owner")
	flags.Int64VarP(&opts.page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.pageSize, "page-size", "", 10, "Size of per page")
	flags.BoolVarP(&opts.public, "public", "", true, "Project is public or private")
	flags.StringVarP(&opts.q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.sort, "sort", "", "", "Sort the resource list in ascending or descending order")
	flags.BoolVarP(&opts.withDetail, "with-detail", "", true, "Bool value indicating whether return detailed information of the project")

	return cmd
}

func RunListProject(opts listProjectOptions) (*project.ListProjectsOK, error) {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.ListProjects(ctx, &project.ListProjectsParams{Name: &opts.name, Owner: &opts.owner, Page: &opts.page, PageSize: &opts.pageSize, Public: &opts.public, Q: &opts.q, Sort: &opts.sort, WithDetail: &opts.withDetail})
	if err != nil {
		return nil, err
	}
	return response, nil
}
