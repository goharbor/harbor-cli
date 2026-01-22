package gc

import (
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/views/gc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ListGCCommand() *cobra.Command {
	var opts api.ListFlags

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List GC history",
		Run: func(cmd *cobra.Command, args []string) {

			history, err := api.GetGCHistory(opts)
			if err != nil {
				logrus.Fatalf("Failed to get GC history: %v", err)
			}

			gc.ListGC(history)
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&opts.Page, "page", "p", 1, "Page number")
	flags.Int64VarP(&opts.PageSize, "page-size", "s", 10, "Size of per page")
	flags.StringVarP(&opts.Sort, "sort", "", "", "Sort the resource list in ascending or descending order")
	flags.StringVarP(&opts.Q, "query", "q", "", "Query string to query resources")

	return cmd
}
