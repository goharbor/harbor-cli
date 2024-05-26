package repository

import (
	"context"
	"github.com/goharbor/harbor-cli/pkg/views/repository/list"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type listRepositoryOptions struct {
	page     int64
	pageSize int64
	q        string
	sort     string
}

func ListRepositoryCommand() *cobra.Command {
	var opts listRepositoryOptions
	var projectName string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list repositories within a project",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName = utils.GetProjectNameFromUser()
			}
			repos, err := runListRepository(projectName, opts)
			if err != nil {
				log.Fatalf("failed to get repositories list: %v", err)
			}
			list.ListRepositories(repos.Payload)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.pageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.sort, "sort", "", "", "Sort the resource list in ascending or descending order")

	return cmd
}

func runListRepository(projectName string, opts listRepositoryOptions) (*repository.ListRepositoriesOK, error) {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Repository.ListRepositories(ctx, &repository.ListRepositoriesParams{ProjectName: projectName, Page: &opts.page, PageSize: &opts.pageSize, Q: &opts.q, Sort: &opts.sort})
	if err != nil {
		return nil, err
	}
	return response, nil
}
