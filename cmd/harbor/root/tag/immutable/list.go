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
	"fmt"

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
		Use:   "list [PROJECT_NAME]",
		Short: "Display all immutable tag rules for a project",
		Long: `Retrieve and display a list of immutable tag rules configured for a specified project in Harbor. 
Immutable tag rules prevent specific tags from being deleted or overwritten, ensuring better security and compliance.
You can specify the project name as an argument or, if omitted, you will be prompted to select one interactively.`,
		Example: `  
  # List immutable tag rules for a specific project  
  harbor tag immutable list my-project  

  # List immutable tag rules interactively (if no project name is provided)  
  harbor tag immutable list  
  `,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var resp immutable.ListImmuRulesOK

			if len(args) > 0 {
				projectName := args[0]
				resp, err = api.ListImmutable(projectName)
			} else {
				projectName, err := prompt.GetProjectNameFromUser()
				if err != nil {
					log.Errorf("failed to get project name: %v", utils.ParseHarborError(err))
				}
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
			if len(resp.Payload) == 0 {
				fmt.Println("No immutable tag rules found.")
				return
			}
			list.ListImmuRules(resp.Payload)
		},
	}
	return cmd
}
