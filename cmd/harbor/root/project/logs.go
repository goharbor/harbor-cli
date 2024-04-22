package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LogsProjectCommmand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "get project logs",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := runLogsProject(args[0]); err != nil {
				log.Errorf("failed to get project logs: %v", err)
			}
		},
	}

	return cmd
}

func runLogsProject(projectName string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	response, err := client.Project.GetLogs(ctx, &project.GetLogsParams{
		ProjectName: projectName,
	})

	if err != nil {
		return err
	}

	log.Infof("Logs of project %s: %s", projectName, response)

	return nil
}
