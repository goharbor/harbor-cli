package tag

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/tag/retention"
	"github.com/spf13/cobra"
)

func TagCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags in Harbor registry",
		Long:  "Manage tags in the Harbor registry, including creating, listing, and deleting retention policies.",
	}
	cmd.AddCommand(retention.Retention())
	return cmd
}
