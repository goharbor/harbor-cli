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
