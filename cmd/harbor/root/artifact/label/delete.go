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
package label

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DelLabelArtifactCommmand delete label command to artifact
func DelLabelArtifactCommmand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "del label to an artifact",
		Long:    `del label to artifact`,
		Example: `harbor artifact label del <project>/<repository>/<reference> <label name>`,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var projectName, repoName, reference string

			if len(args) > 0 {
				projectName, repoName, reference = utils.ParseProjectRepoReference(args[0])
			}

			labels, err := api.ListLabel()
			if err != nil {
				log.Errorf("failed to list label: %v", err)
				return
			}

			var label *models.Label
			for _, currentLabel := range labels.GetPayload() {
				if currentLabel.Name == args[1] {
					label = currentLabel
				}
			}

			_, err = api.RemoveLabelArtifact(projectName, repoName, reference, label)
			if err != nil {
				log.Errorf("failed to remove label on artifact: %v", err)
				return
			}

			log.Infof("Label %s removed on artifact %s.", args[1], args[0])
		},
	}

	return cmd
}
