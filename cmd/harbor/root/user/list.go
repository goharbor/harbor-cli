package user

import (
	"context"
	log "github.com/sirupsen/logrus"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/user/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type listUserOptions struct {
	pageSize int64
	page     int64
	q        string
	sort     string
}

func UserListCmd() *cobra.Command {
	var opts listUserOptions

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list users",
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		Run: func(cmd *cobra.Command, args []string) {
			response, err := runListUsers(opts)
			if err != nil {
				log.Fatalf("failed to get users list: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(response.Payload)
			} else {
				list.ListUsers(response.Payload)
			}
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.pageSize, "page-size", "", 10, "Size of per page")
	flags.Int64VarP(&opts.page, "page", "", 1, "Page number")
	flags.StringVarP(&opts.q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.sort, "sort", "", "", "Sort the resource list in ascending or descending order")

	return cmd
}

func runListUsers(opts listUserOptions) (*user.ListUsersOK, error) {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.User.ListUsers(ctx, &user.ListUsersParams{Page: &opts.page, PageSize: &opts.pageSize, Q: &opts.q, Sort: &opts.sort})
	if err != nil {
		return nil, err
	}
	return response, nil
}
