package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type candidateScannerOptions struct {
	projectNameOrId string
	pageSize        int64
	page            int64
	q               string
	sort            string
}

func GetCandidateScanner() *cobra.Command {
	var opts candidateScannerOptions

	cmd := &cobra.Command{
		Use:   "candidate [Name|ID]",
		Short: "Get scanner registration candidates for configuring project level scanner",
		Long:  "Retrieve the system configured scanner registrations as candidates of setting project level scanner.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				opts.projectNameOrId = args[0]
			} else {
				projectName := utils.GetProjectNameFromUser()
				opts.projectNameOrId = projectName
			}

			response, err := runGetCandidateScanner(opts)
			if err != nil {
				log.Fatalf("failed to get candidate scanner: %v", err)
			}
			utils.PrintPayloadInJSONFormat(response)
		},
	}

	flags := cmd.Flags()

	flags.Int64VarP(&opts.pageSize, "page-size", "", 10, "Size of per page")
	flags.Int64VarP(&opts.page, "page", "", 1, "Page number")
	flags.StringVarP(&opts.q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.sort, "sort", "", "", "Sort the resource list in ascending or descending order. e.g. sort by field1 in ascending order and field2 in descending order with 'sort=field1,-field2'")

	return cmd

}

func runGetCandidateScanner(opts candidateScannerOptions) (*project.ListScannerCandidatesOfProjectOK, error) {

	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.ListScannerCandidatesOfProject(ctx, &project.ListScannerCandidatesOfProjectParams{ProjectNameOrID: opts.projectNameOrId, PageSize: &opts.pageSize, Page: &opts.page, Q: &opts.q, Sort: &opts.sort})
	if err != nil {
		return nil, err
	}

	return response, nil
}
