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
package purge

import (
	"github.com/spf13/cobra"
)

func ListPurgeCommand() *cobra.Command {
	var (
		page     int64
		pageSize int64
		query    string
		sort     string
		purgeID  int64
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get purge job execution results",
		Args:  cobra.ExactArgs(0),
	}

	flags := cmd.Flags()
	flags.Int64VarP(&page, "page", "", 1, "Page number")
	flags.Int64VarP(&pageSize, "page-size", "", 20, "Size per page (max 100)")
	flags.StringVarP(&query, "query", "q", "", "Query filter to search purge job results")
	flags.StringVarP(&sort, "sort", "", "", "Sort the purge job results in ascending or descending order")
	flags.Int64VarP(&purgeID, "purge-id", "p", 0, "ID of the purge job to filter results")

	return cmd
}
