package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DeleteProjectCommand creates a new `harbor delete project` command
func DeleteProjectCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete project by name or id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			credentialName := viper.GetString("current-credential-name")
			client := utils.GetClientByCredentialName(credentialName)
			ctx := context.Background()
			if len(args) > 0 {
				err = RunDeleteProject(args[0], ctx, client.Project)
			} else {
				projectName := utils.GetProjectNameFromUser()
				err = RunDeleteProject(projectName, ctx, client.Project)
			}
			if err != nil {
				log.Errorf("failed to delete project: %v", err)
			}
		},
	}

	return cmd
}

func RunDeleteProject(projectName string, ctx context.Context, projectInterface ProjectInterface) error {

	_, err := projectInterface.DeleteProject(ctx, &project.DeleteProjectParams{ProjectNameOrID: projectName})

	if err != nil {
		return err
	}

	log.Info("project deleted successfully")
	return nil
}
