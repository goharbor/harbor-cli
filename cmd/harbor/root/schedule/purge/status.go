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

func StatusCommand() *cobra.Command {
	var (
		purgeID int64
	)

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the status of a purge job by ID",
		Long:  "Show the status of a purge job by ID in Harbor",
		Args:  cobra.MaximumNArgs(1),
	}

	cmd.Flags().Int64VarP(&purgeID, "purge-id", "p", 0, "ID of the purge job to stop")

	return cmd
}
