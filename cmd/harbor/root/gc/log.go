package gc

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GetGCLogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log [gc_id]",
		Short: "Get GC job log",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				logrus.Fatalf("Invalid GC ID: %v", err)
			}

			logData, err := api.GetGCJobLog(id)
			if err != nil {
				logrus.Fatalf("Failed to get GC log: %v", err)
			}

			fmt.Println(logData)
		},
	}
	return cmd
}
