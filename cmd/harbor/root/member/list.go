package member

import (
	"context"
	"os"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/member"
	"github.com/goharbor/harbor-cli/pkg/utils"
	list "github.com/goharbor/harbor-cli/pkg/views/member/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type listMemberOptions struct {
	projectNameOrID string
	page            int64
	pageSize        int64
	entityName      string
	withDetail      bool
}

// ListMemberCommand creates a new `harbor member list` command
func ListMemberCommand() *cobra.Command {
	var opts listMemberOptions

	cmd := &cobra.Command{
		Use:     "list [projectName or ID]",
		Short:   "list members of a project by projectName Or ID",
		Args:    cobra.MaximumNArgs(1),
		Example: "harbor member list my-project",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				opts.projectNameOrID = args[0]
			} else {
				opts.projectNameOrID = utils.GetProjectNameFromUser()
				if opts.projectNameOrID == "" {
					os.Exit(1)
				}
			}

			members, err := RunListMember(opts)
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
	flags.Int64VarP(&opts.page, "page", "", 1, "Page number")
	flags.Int64VarP(&opts.pageSize, "page-size", "", 10, "Size of per page")
	flags.StringVarP(&opts.entityName, "name", "n", "", "Member Name to search")

	return cmd
}

func RunListMember(opts listMemberOptions) (*member.ListProjectMembersOK, error) {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Member.ListProjectMembers(
		ctx,
		&member.ListProjectMembersParams{
			ProjectNameOrID: opts.projectNameOrID,
			Entityname:      &opts.entityName,
			Page:            &opts.page,
			PageSize:        &opts.pageSize,
		},
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}
