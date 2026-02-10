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

func LogsCommand() *cobra.Command {
	var (
		purgeID int64
	)

	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get purge job logs by ID",
		Long:  "Get purge job logs filtered by specific purge ID.",
		Args:  cobra.MaximumNArgs(1),
	}

	cmd.Flags().Int64VarP(&purgeID, "purge-id", "p", 0, "ID of the purge job to retrieve logs for")

	return cmd
}
