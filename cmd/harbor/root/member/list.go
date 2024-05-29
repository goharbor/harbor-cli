package member

import (
	"os"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/member/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListMemberCommand creates a new `harbor member list` command
func ListMemberCommand() *cobra.Command {
	var opts api.ListMemberOptions

	cmd := &cobra.Command{
		Use:     "list [projectName or ID]",
		Short:   "list members in a project",
		Long:    "list members in a project by projectName Or ID",
		Example: "  harbor member list my-project",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				opts.ProjectNameOrID = args[0]
			} else {
				opts.ProjectNameOrID = prompt.GetProjectNameFromUser()
				if opts.ProjectNameOrID == "" {
					os.Exit(1)
				}
			}

			members, err := api.ListMember(opts)
			if err != nil {
				log.Fatalf("failed to get members list: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag == "json" {
				utils.PrintPayloadInJSONFormat(members)
				return
			}

			if FormatFlag == "wide" {
				list.ListMembers(members.Payload, true)
			} else {
				list.ListMembers(members.Payload, false)
			}
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.EntityName, "name", "n", "", "Member Name to search")

	return cmd
}
