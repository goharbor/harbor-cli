package config

import (
	"fmt"
	"strings"

	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func AddConfigCommand() *cobra.Command {
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

				err := api.AddConfig(isID, projectNameOrID, metadata)
				if err != nil {
					log.Errorf("failed to add metadata: %v", err)
				}
			}

		},
	}

	return cmd
}
