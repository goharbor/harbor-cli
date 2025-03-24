package tag

import (
	"github.com/goharbor/harbor-cli/cmd/harbor/root/tag/retention"
	"github.com/spf13/cobra"
)

var TagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
	Long:  `Manage tags in the Harbor project.`,
}

func init() {
	TagCmd.AddCommand(retention.Retention())
}
