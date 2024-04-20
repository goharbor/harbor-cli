package project

import (
	"context"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project"
	"github.com/goharbor/harbor-cli/pkg/constants"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DeleteProjectCommand creates a new `harbor delete project` command
func DeleteProjectCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete [NAME|ID]",
		Short: "delete project by name or id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)

			if len(args) > 0 {
				err = runDeleteProject(args[0], credentialName)
			} else {
				projectName := utils.GetProjectNameFromUser(credentialName)
				err = runDeleteProject(projectName, credentialName)
			}
			if err != nil {
				log.Errorf("failed to delete project: %v", err)
			}
		},
	}

	return cmd
}

func runDeleteProject(projectName string, credentialName string) error {
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Project.DeleteProject(ctx, &project.DeleteProjectParams{ProjectNameOrID: projectName})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}
