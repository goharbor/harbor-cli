package project

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetProjectCommand creates a new `harbor get project` command
func ViewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "view [NAME|ID]",
		Short: "get project by name or id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var projectName string
			var project *project.GetProjectOK

			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName = prompt.GetProjectNameFromUser()
			}

			project, err = api.GetProject(projectName)

			if err != nil {
				log.Errorf("failed to get project: %v", err)
				return
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(project, FormatFlag)
				if err != nil {
					log.Error(err)
					return
				}
			} else {
				view.ViewProjects(project.Payload)
			}

		},
	}

	return cmd
}
