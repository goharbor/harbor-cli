package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type getProjectSummaryOptions struct {
	projectNameOrID string
}

func SummaryCommand() *cobra.Command {
	var opts getProjectSummaryOptions

	cmd := &cobra.Command{
		Use:   "summary [NAME|ID]",
		Short: "Get project summary by name or id",
		Long:  "Get project summary by name or id. If the project name or id is not provided, the command will prompt the user to enter the project name.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) > 0 {
				opts.projectNameOrID = args[0]
			} else {
				projectName := utils.GetProjectNameFromUser()
				opts.projectNameOrID = projectName
			}

			response, err := runGetProjectSummary(opts)
			if err != nil {
				log.Fatalf("failed to get project summary: %v", err)
			}
			utils.PrintPayloadInJSONFormat(response.GetPayload())

		},
	}

	return cmd
}

func runGetProjectSummary(opts getProjectSummaryOptions) (*project.GetProjectSummaryOK, error) {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.GetProjectSummary(ctx, &project.GetProjectSummaryParams{ProjectNameOrID: opts.projectNameOrID})

	if err != nil {
		return nil, err
	}

	return response, nil
}
