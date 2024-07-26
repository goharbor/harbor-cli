package user

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/user/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UserListCmd() *cobra.Command {
	var formatFlag string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list users",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			response, err := api.ListUsers()
			if err != nil {
				log.Errorf("failed to list users: %v", err)
				return
			}
			if formatFlag != "" {
				if formatFlag == "json" {
					utils.PrintPayloadInJSONFormat(response)
				} else if formatFlag == "yaml" {
					utils.PrintPayloadInYAMLFormat(response)
				} else {
					log.Errorf("invalid output format: %s", formatFlag)
				}
			} else {
				list.ListUsers(response.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&formatFlag, "output-format", "o", "", "Output format. One of: json|yaml")

	return cmd
}
