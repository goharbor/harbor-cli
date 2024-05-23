package replication

import (
	"github.com/spf13/cobra"
)

func Replication() *cobra.Command {
	// replicationCmd represents the replication command.
	var replicationCmd = &cobra.Command{
		Use:     "replication",
		Aliases: []string{"repl"},
		Short:   "",
		Long:    ``,
	}
	replicationCmd.AddCommand()

	return replicationCmd
}
