// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package user

import (
	"fmt"
	"io"
	"os"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/user/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetUsers(opts api.ListFlags) ([]*models.UserResp, error) {
	var allUsers []*models.UserResp

	if opts.PageSize > 100 {
		return nil, fmt.Errorf("page size should be less than or equal to 100")
	}

	if opts.PageSize == 0 {
		opts.PageSize = 100
		opts.Page = 1

		for {
			response, err := api.ListUsers(opts)
			if err != nil {
				if isUnauthorizedError(err) {
					return nil, fmt.Errorf("Permission denied: Admin privileges are required to execute this command.")
				}
				return nil, fmt.Errorf("failed to list users: %v", err)
			}

			allUsers = append(allUsers, response.Payload...)

			if len(response.Payload) < int(opts.PageSize) {
				break
			}
			opts.Page++
		}
	} else {
		response, err := api.ListUsers(opts)
		if err != nil {
			if isUnauthorizedError(err) {
				return nil, fmt.Errorf("Permission denied: Admin privileges are required to execute this command.")
			}
			return nil, fmt.Errorf("failed to list users: %v", err)
		}
		allUsers = response.Payload
	}
	return allUsers, nil
}
func PrintUsers(w io.Writer, allUsers []*models.UserResp) error {
	if len(allUsers) == 0 {
		log.Info("No users found")
		return nil
	}
	formatFlag := viper.GetString("output-format")
	if formatFlag != "" {
		err := utils.FPrintFormat(w, allUsers, formatFlag)
		if err != nil {
			log.Error(err)
		}
	} else {
		if err := list.ListUsers(w, allUsers); err != nil {
			return err
		}
	}
	return nil
}
func UserListCmd() *cobra.Command {
	var opts api.ListFlags
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List users",
		Args:    cobra.ExactArgs(0),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			allUsers, err := GetUsers(opts)
			if err != nil {
				return err
			}
			return PrintUsers(os.Stdout, allUsers)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "p", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "n", 0, "Size of per page (0 to fetch all)")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "s", "", "Sort the resource list in ascending or descending order")

	return cmd
}
