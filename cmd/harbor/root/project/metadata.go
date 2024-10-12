package project

import (
	"fmt"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

var (
	isID bool
)

func ProjectMetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadata",
		Short: "Manage project metadata",
	}
	cmd.AddCommand(
		AddMetadataCommand(),
		DeleteMetadataCommand(),
		ViewMetadataCommand(),
		UpdateMetadataCommand(),
		ListMetadataCommand(),
	)

	flags := cmd.Flags()
	flags.BoolVarP(&isID, "id", "", false, "Use project ID instead of name")

	return cmd
}

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

				err := api.AddMetadata(isID, projectNameOrID, metadata)
				if err != nil {
					log.Errorf("failed to add metadata: %v", err)
				}
			}

		},
	}

	return cmd
}

func DeleteMetadataCommand() *cobra.Command {
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

				err := api.DeleteMetadata(isID, projectNameOrID, metaNames)
				if err != nil {
					log.Errorf("failed to delete metadata: %v", err)
				}
			}

		},
	}

	return cmd
}

func ListMetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list [NAME|ID]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name")
			} else {
				projectNameOrID := args[0]

				response, err := api.ListMetadata(isID, projectNameOrID)
				if err != nil {
					log.Errorf("failed to view metadata: %v", err)
				} else {
					utils.PrintPayloadInJSONFormat(response.Payload)
				}
			}

		},
	}

	return cmd
}

func UpdateMetadataCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update [NAME|ID] [KEY] ...[KEY]:[VALUE]",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println("Please provide project name, the meta name and metadata")
			} else if len(args) == 1 {
				fmt.Println("Please provide the meta name and metadata")
			} else if len(args) == 2 {
				fmt.Println("Please provide metadata")
			} else {
				projectNameOrID := args[0]
				metaName := args[1]
				metadata := make(map[string]string)
				for i := 2; i < len(args); i++ {
					keyValue := args[i]
					keyValueArray := strings.Split(keyValue, ":")
					if len(keyValueArray) == 2 {
						metadata[keyValueArray[0]] = keyValueArray[1]
					} else {
						fmt.Println("Please provide metadata in the format key:value")
						return
					}
				}

				err := api.UpdateMetadata(isID, projectNameOrID, metaName, metadata)
				if err != nil {
					log.Errorf("failed to view metadata: %v", err)
				}
			}

		},
	}

	return cmd
}

func ViewMetadataCommand() *cobra.Command {
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

				err := api.ViewMetadata(isID, projectNameOrID, metaName)
				if err != nil {
					log.Errorf("failed to view metadata: %v", err)
				}
			}

		},
	}

	return cmd
}
