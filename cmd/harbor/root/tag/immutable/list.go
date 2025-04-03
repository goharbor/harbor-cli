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
package immutable

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/immutable"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/immutable/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListImmutableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all immutable tag rule of project",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var resp immutable.ListImmuRulesOK

			if len(args) > 0 {
				projectName := args[0]
				resp, err = api.ListImmutable(projectName)
			} else {
				projectName := prompt.GetProjectNameFromUser()
				resp, err = api.ListImmutable(projectName)
			}

			if err != nil {
				log.Errorf("failed to list immutablility rule: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(resp)
				return
			}

			list.ListImmuRules(resp.Payload)
		},
	}
	return cmd
}
