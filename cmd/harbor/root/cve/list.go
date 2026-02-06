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
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/cveallowlist/list"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List system level allowlist of cve",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cve, err := api.ListSystemCve()
			if err != nil {
				return fmt.Errorf("failed to get system cve list: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				err = utils.PrintFormat(cve, FormatFlag)
				if err != nil {
					return fmt.Errorf("failed to print cve list: %v", err)
				}
			} else {
				list.ListSystemCve(cve.Payload)
			}
			return nil
		},
	}

	return cmd
}
