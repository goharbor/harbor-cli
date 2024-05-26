package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/utils"
	auditLog "github.com/goharbor/harbor-cli/pkg/views/project/logs"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type logsProjectOptions struct {
	page     int64
	pageSize int64
	q        string
	sort     string
}

func LogsProjectCommmand() *cobra.Command {
	var opts logsProjectOptions

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "get project logs",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var resp *project.GetLogsOK
			if len(args) > 0 {
				resp, err = runLogsProject(args[0], opts)
			} else {
				projectName := utils.GetProjectNameFromUser()
				resp, err = runLogsProject(projectName, opts)
			}

			if err != nil {
				log.Fatalf("failed to get project logs: %v", err)
			}
			auditLog.LogsProject(resp.Payload)

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(resp)
				return
			}

		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.pageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.sort, "sort", "", "", "Sort the resource list in ascending or descending order")

	return cmd
}

func runLogsProject(projectName string, opts logsProjectOptions) (*project.GetLogsOK, error) {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.GetLogs(ctx, &project.GetLogsParams{
		ProjectName: projectName,
		Page:        &opts.page,
		PageSize:    &opts.pageSize,
		Q:           &opts.q,
		Sort:        &opts.sort,
		Context:     ctx,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
