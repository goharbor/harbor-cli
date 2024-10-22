package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type projectLevelScannerOptions struct {
	projectNameOrID string
}

func ScannerCommand() *cobra.Command {
	var opts projectLevelScannerOptions

	cmd := &cobra.Command{
		Use:   "scanner [NAME|ID]",
		Short: "Get the scanner registration of a project.",
		Long:  "Get the scanner registration of a project by name or id. If no scanner registration is configured for the specified project, the system default scanner registration will be returned",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				opts.projectNameOrID = args[0]
			} else {
				projectName := utils.GetProjectNameFromUser()
				opts.projectNameOrID = projectName
			}

			response, err := runGetProjectScanner(opts)
			if err != nil {
				log.Fatalf("failed to get project scanner: %v", err)
			}
			utils.PrintPayloadInJSONFormat(response.GetPayload())

		},
	}
	return cmd
}

func runGetProjectScanner(opts projectLevelScannerOptions) (*project.GetScannerOfProjectOK, error) {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.GetScannerOfProject(ctx, &project.GetScannerOfProjectParams{ProjectNameOrID: opts.projectNameOrID})

	if err != nil {
		return nil, err
	}

	return response, nil
}
