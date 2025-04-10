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
package webhook

import "github.com/spf13/cobra"

func Webhook() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "webhook",
		Short:   "Manage webhooks",
		Long:    `Manage webhooks in Harbor Repository`,
		Example: `  harbor webhook list`,
	}
	cmd.AddCommand(
		CreateWebhookCmd(),
		ListWebhookCommand(),
		DeleteWebhookCmd(),
		EditWebhookCmd(),
	)
	return cmd
}
