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

type viewMetadataOptions struct {
	isID bool
}

func ViewMetadataCommand() *cobra.Command {
	var opts viewMetadataOptions

	cmd := &cobra.Command{
		Use:   "view",
		Short: "view [NAME|ID] [KEY]",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name or id and the meta name")
			} else if len(args) == 1 {
				fmt.Println("Please provide the meta name")
			} else {
				projectNameOrID := args[0]
				metaName := args[1]

				err := viewMetadata(opts, projectNameOrID, metaName)
				if err != nil {
					log.Errorf("failed to view metadata: %v", err)
				}
			}

		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.isID, "id", "", false, "Use project ID instead of name")

	return cmd
}

func viewMetadata(opts viewMetadataOptions, projectNameOrID string, metaName string) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()

	isName := !opts.isID
	response, err := client.ProjectMetadata.GetProjectMetadata(ctx, &project_metadata.GetProjectMetadataParams{MetaName: metaName, ProjectNameOrID: projectNameOrID, XIsResourceName: &isName})
	if err != nil {
		return err
	}
	if response != nil {
		log.Info("Metadata: ", response.Payload)
	}

	return nil
}
