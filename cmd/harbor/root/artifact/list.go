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
package artifact

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/artifact"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	artifactViews "github.com/goharbor/harbor-cli/pkg/views/artifact/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListArtifactCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List container artifacts (images, charts, etc.) in a Harbor repository with metadata",
		Long: `List all artifacts (e.g., container images, charts) within a given Harbor repository. 
Supports optional project/repository input in the form <project>/<repository>. 
Displays key artifact metadata including tags, digest, type, size, vulnerability count, and push time.

Examples:
  harbor-cli artifact list                # Interactive prompt for project and repository
  harbor-cli artifact list library/nginx # Directly list artifacts in the nginx repo under 'library' project

Supports pagination, search queries, and sorting using flags.`,

		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.PageSize < 0 {
				return fmt.Errorf("page size must be greater than or equal to 0")
			}

			if opts.PageSize > 100 {
				return fmt.Errorf("page size should be less than or equal to 100")
			}
			var err error
			var artifacts artifact.ListArtifactsOK
			var projectName, repoName string

			if len(args) > 0 {
				projectName, repoName, err = utils.ParseProjectRepo(args[0])
				if err != nil {
					return fmt.Errorf("failed to parse project/repo: %v", err)
				}
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
				repoName = prompt.GetRepoNameFromUser(projectName)
			}

			artifacts, err = api.ListArtifact(projectName, repoName, opts)

			if err != nil {
				return fmt.Errorf("failed to list artifacts: %v", err)
			}

			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(artifacts, FormatFlag)
				if err != nil {
					return err
				}
			} else {
				artifactViews.ListArtifacts(artifacts.Payload)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "p", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "n", 10, "Size of per page")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")
	flags.StringVarP(&opts.Sort, "sort", "s", "", "Sort the resource list in ascending or descending order")

	return cmd
}
