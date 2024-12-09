package project

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/project/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SearchProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "search project based on their names",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			projects, err := api.SearchProject(args[0])
			if err != nil {
				log.Fatalf("failed to get projects: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(projects, FormatFlag)
				if err != nil {
					log.Error(err)
				}
			} else {
				list.SearchProjects(projects.Payload.Project)
			}

		},
	}
	return cmd
}
