package member

import (
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	view "github.com/goharbor/harbor-cli/pkg/views/member/view"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewGetRegistryCommand creates a new `harbor get registry` command
func ViewMemberCommand() *cobra.Command {
	var opts api.GetMemberOptions

	cmd := &cobra.Command{
		Use:     "view [ProjectName Or ID] [member ID]",
		Short:   "get project member by ID",
		Long:    "get member details by MemberID",
		Example: "  harbor project member view my-project [memberID]",
		Args:    cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				opts.ProjectNameOrID = args[0]
			} else if len(args) == 2 {
				opts.ID, _ = strconv.ParseInt(args[1], 0, 64)
				opts.ProjectNameOrID = args[0]
			} else if opts.ProjectNameOrID == "" || opts.ID == 0 {
				if opts.ProjectNameOrID == "" {
					opts.ProjectNameOrID = prompt.GetProjectNameFromUser()
				}
				if opts.ID == 0 {
					opts.ID = prompt.GetMemberIDFromUser(opts.ProjectNameOrID)
				}
			}

			member, err := api.GetMember(opts)
			if err != nil {
				log.Fatalf("failed to get members list: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			VerboseFlag := viper.GetBool("verbose")
			if FormatFlag == "json" {
				utils.PrintPayloadInJSONFormat(member)
				return
			} else if VerboseFlag {
				view.ViewMember(member.Payload, true)
			} else {
				view.ViewMember(member.Payload, false)
			}
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.ID, "id", "", 0, "Member ID")
	flags.StringVarP(&opts.ProjectNameOrID, "projectname", "p", "", "Project Name")
	return cmd
}
