package repository

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/repository"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/repository/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListRepositoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list repositories within a project",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var repos repository.ListRepositoriesOK
			var projectName string

			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName = prompt.GetProjectNameFromUser()
			}

			repos, err = api.ListRepository(projectName)

			if err != nil {
				log.Errorf("failed to list repositories: %v", err)
			}
			
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				if FormatFlag == "json" {
					utils.PrintPayloadInJSONFormat(repos)
					return
				}
				if FormatFlag == "yaml" {
					utils.PrintPayloadInYAMLFormat(repos)
					return
				}
				log.Errorf("Unable to output in the specified '%s' format", FormatFlag)
				return
			}

			list.ListRepositories(repos.Payload)

		},
	}

	return cmd
}
