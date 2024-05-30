package metadata

import (
	"context"
	"fmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/project_metadata"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type deleteMetadataOptions struct {
	isID bool
}

func DeleteMetadataCommand() *cobra.Command {
	var opts deleteMetadataOptions

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete [NAME|ID] ...[KEY]",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name or id and the meta names to delete")
			} else if len(args) == 1 {
				fmt.Println("Please provide the meta names")
			} else {
				projectNameOrID := args[0]
				metaNames := make([]string, 0)
				for i := 1; i < len(args); i++ {
					metaNames = append(metaNames, args[i])
				}

				err := deleteMetadata(opts, projectNameOrID, metaNames)
				if err != nil {
					log.Errorf("failed to delete metadata: %v", err)
				}
			}

		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.isID, "id", "", false, "Use project ID instead of name")

	return cmd
}

func deleteMetadata(opts deleteMetadataOptions, projectNameOrID string, metaName []string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	isName := !opts.isID
	for _, meta := range metaName {
		response, err := client.ProjectMetadata.DeleteProjectMetadata(ctx, &project_metadata.DeleteProjectMetadataParams{MetaName: meta, ProjectNameOrID: projectNameOrID, XIsResourceName: &isName})
		if err != nil {
			return err
		}
		if response != nil {
			log.Info("Metadata %v deleted successfully", meta)
		}
	}

	return nil
}
