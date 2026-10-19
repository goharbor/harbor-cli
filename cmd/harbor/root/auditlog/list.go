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
package auditlog

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	auditlogView "github.com/goharbor/harbor-cli/pkg/views/auditlog"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListAuditLogsCommand() *cobra.Command {
	var (
		opts      api.ListFlags
		username  string
		operation string
		resource  string
		query     string
	)

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List audit logs",
		Long:    "List system-wide audit logs in Harbor",
		Example: `harbor auditlog list --page 1 --page-size 10 --username admin`,
		Run: func(cmd *cobra.Command, args []string) {
			var matchFilters []string
			if username != "" {
				matchFilters = append(matchFilters, "username="+username)
			}
			if operation != "" {
				matchFilters = append(matchFilters, "operation="+operation)
			}
			if resource != "" {
				matchFilters = append(matchFilters, "resource="+resource)
			}

			if len(matchFilters) > 0 {
				q, err := utils.BuildQueryParam(nil, matchFilters, nil, []string{"username", "operation", "resource"})
				if err != nil {
					logrus.Fatalf("failed to build query parameters: %v", err)
				}
				opts.Q = q
			} else if query != "" {
				opts.Q = query
			}

			response, err := api.AuditLogs(opts)
			if err != nil {
				logrus.Fatalf("failed to fetch audit logs: %v", err)
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(response.Payload, formatFlag)
				if err != nil {
					logrus.Fatalf("failed to print format: %v", err)
				}
				return
			}

			auditlogView.RenderAuditLogsExt(response.Payload)
		},
	}

	cmd.Flags().Int64VarP(&opts.Page, "page", "", 1, "Page number")
	cmd.Flags().Int64VarP(&opts.PageSize, "page-size", "", 10, "Size of per page")
	cmd.Flags().StringVarP(&query, "query", "q", "", "Query string for filtering logs")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Filter logs by username")
	cmd.Flags().StringVarP(&operation, "operation", "", "", "Filter logs by operation type (e.g. create, delete, pull)")
	cmd.Flags().StringVarP(&resource, "resource", "r", "", "Filter logs by target resource")

	return cmd
}
