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
package labels

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteLabelCommand() *cobra.Command {
	var opts models.Label
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete label",
		Example: "harbor label delete [labelname]",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			deleteView := &api.ListFlags{
				Scope: opts.Scope,
			}

			if len(args) > 0 {
				labelId, _ := api.GetLabelIdByName(args[0])
				err = api.DeleteLabel(labelId)
			} else {
				labelId := prompt.GetLabelIdFromUser(*deleteView)
				err = api.DeleteLabel(labelId)
			}
			if err != nil {
				log.Errorf("failed to delete label: %v", err)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.Scope, "scope", "s", "g", "default(global).'p' for project labels.Query scope of the label")

	return cmd
}
