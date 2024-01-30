package robot

import (
	"context"

	"github.com/akshatdalton/harbor-cli/cmd/constants"
	"github.com/akshatdalton/harbor-cli/cmd/utils"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/robot"
	"github.com/spf13/cobra"
)

type listRobotOptions struct {
	page     int64
	pageSize int64
	q        string
	sort     string
}

// NewListRobotCommand creates a new `harbor list robot` command
func NewListRobotCommand() *cobra.Command {
	var opts listRobotOptions

	cmd := &cobra.Command{
		Use:   "robot",
		Short: "list robot",
		RunE: func(cmd *cobra.Command, args []string) error {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return runListRobot(opts, credentialName)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.pageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.sort, "sort", "", "", "Sort the resource list in ascending or descending order")

	return cmd
}

func runListRobot(opts listRobotOptions, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Robot.ListRobot(ctx, &robot.ListRobotParams{Page: &opts.page, PageSize: &opts.pageSize, Q: &opts.q, Sort: &opts.sort})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.GetPayload())
	return nil
}
