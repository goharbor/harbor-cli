package repository

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/repository/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ListRepositoryCommand() *cobra.Command {
	var formatFlag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list [PROJECT]",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var response repository.ListRepositoriesOK

			if len(args) > 0 {
				response, err = api.ListRepository(args[0])
			} else {
				projectName := prompt.GetProjectNameFromUser()
				response, err = api.ListRepository(projectName)
			}
			if err != nil {
				log.Errorf("failed to list repositories: %v", err)
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
				list.ListRepositories(response.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&formatFlag, "output-format", "o", "", "Output format. One of: json|yaml")

	return cmd
}
