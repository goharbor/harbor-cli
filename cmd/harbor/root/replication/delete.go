package replication

import (
	"fmt"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
	"strconv"
)

func DeleteReplicationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"repl"},
		Short:   "delete replication policies",
		Long:    `delete replication policies`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			replicationPolicyID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("invalid project id: ", utils.ParseHarborErrorMsg(err))
				return
			}

			err = api.DeleteReplication(replicationPolicyID)
			if err != nil {
				fmt.Printf("failed to delete replication policy %d: %v\n", replicationPolicyID, utils.ParseHarborErrorMsg(err))
				return
			} else {
				fmt.Printf("deleted replication policy %d\n", replicationPolicyID)
			}
		},
	}

	return cmd
}
