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
package cve

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/cveallowlist/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list system level allowlist of cve",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cve, err := api.ListSystemCve()
			if err != nil {
				log.Fatalf("failed to get system cve list: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(cve, FormatFlag)
				if err != nil {
					log.Fatalf("failed to print cve list: %v", err)
					return
				}
			} else {
				list.ListSystemCve(cve.Payload)
			}
		},
	}

	return cmd
}
