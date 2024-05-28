package metadata

import (
	"context"
	"fmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project_metadata"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func AddMetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add [NAME|ID] ...[KEY]:[VALUE]",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name or id and the metadata")
			} else if len(args) == 1 {
				fmt.Println("Please provide the metadata")
			} else {
				projectNameOrID := args[0]
				metadata := make(map[string]string)
				for i := 1; i < len(args); i++ {
					keyValue := args[i]
					keyValueArray := strings.Split(keyValue, ":")
					if len(keyValueArray) == 2 {
						metadata[keyValueArray[0]] = keyValueArray[1]
					} else {
						fmt.Println("Please provide metadata in the format key:value")
						return
					}
				}

				err := addMetadata(projectNameOrID, metadata)
				if err != nil {
					log.Errorf("failed to add metadata: %v", err)
				}
			}

		},
	}

	return cmd
}

func addMetadata(projectNameOrID string, metadata map[string]string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	response, err := client.ProjectMetadata.AddProjectMetadatas(ctx, &project_metadata.AddProjectMetadatasParams{Metadata: metadata, ProjectNameOrID: projectNameOrID})
	if err != nil {
		return err
	}
	if response != nil {
		log.Info("Metadata added successfully")
	}

	return nil
}
