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
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AuditLogTypesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "types",
		Short:   "List audit log event types",
		Long:    "List available event types for audit logs in Harbor",
		Example: `harbor auditlog types`,
		Run: func(cmd *cobra.Command, args []string) {
			response, err := api.AuditLogEventTypes()
			if err != nil {
				logrus.Fatalf("failed to fetch audit log event types: %v", err)
			}

			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				err = utils.PrintFormat(response.Payload, formatFlag)
				if err != nil {
					logrus.Fatalf("failed to print format: %v", err)
				}
				return
			}

			for _, eventType := range response.Payload {
				fmt.Println(eventType)
			}
		},
	}

	return cmd
}
