package robot

import (
	"log"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListRobotCommand creates a new `harbor robot list` command
func ListRobotCommand() *cobra.Command {
	var (
		query     string
		opts      api.ListFlags
	)

	projectQString := constants.ProjectQString
	cmd := &cobra.Command{
		Use:   "list [projectID]",
		Short: "list robot",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				opts.Q = projectQString + args[0]
			} else {
        projectID := prompt.GetProjectIDFromUser()
				opts.Q = projectQString + strconv.FormatInt(projectID, 10)
			}

			robots, err := api.ListRobot(opts)
			if err != nil {
				log.Fatalf("failed to get robots list: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(robots)
				return
			}

			list.ListRobots(robots.Payload)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&query, "query", "q", "", "Query string to query resources")
	flags.StringVarP(
		&opts.Sort,
		"sort",
		"",
		"",
		"Sort the resource list in ascending or descending order",
	)

	return cmd
}
