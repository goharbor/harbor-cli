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
package artifactscan

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func StopScanArtifactCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop",
		Short:   "Stop a scan of an artifact",
		Long:    `Stop a scan of an artifact in Harbor Repository`,
		Example: `harbor artifact scan stop <project>/<repository>/<reference>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			var projectName, repoName, reference string

			if len(args) > 0 {
				projectName, repoName, reference, err = utils.ParseProjectRepoReference(args[0])
				if err != nil {
					return fmt.Errorf("failed to parse project/repo/reference: %v", err)
				}
			} else {
				projectName, err = prompt.GetProjectNameFromUser()
				if err != nil {
					return fmt.Errorf("failed to get project name: %v", utils.ParseHarborErrorMsg(err))
				}
				repoName = prompt.GetRepoNameFromUser(projectName)
				reference = prompt.GetReferenceFromUser(repoName, projectName)
			}

			err = api.StopScanArtifact(projectName, repoName, reference)
			if err != nil {
				return fmt.Errorf("failed to stop scan of artifact: %v", err)
			}

			return nil
		},
	}
	return cmd
}
