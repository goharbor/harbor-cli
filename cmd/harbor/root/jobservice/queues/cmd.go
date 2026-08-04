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
package queues

import "github.com/spf13/cobra"

// QueuesCommand creates the queues subcommand
func QueuesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "queues",
		Short: "Manage job queues (list, stop, pause, resume)",
		Long:  "List job queues and perform actions on them (stop/pause/resume).",
	}

	cmd.AddCommand(ListCommand(), StopCommand(), PauseCommand(), ResumeCommand())

	return cmd
}
