package root

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/statistic"
	"github.com/goharbor/harbor-cli/pkg/utils"
	statisticview "github.com/goharbor/harbor-cli/pkg/views/statistic"
	"github.com/spf13/cobra"
)

func StatisticCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "statistic",
        Short: "Get statistics about Harbor projects and repositories",
        RunE: func(cmd *cobra.Command, args []string) error {
            ctx, client, err := utils.ContextWithClient()
            if err != nil {
                return fmt.Errorf("failed to create client: %w", err)
            }

            params := &statistic.GetStatisticParams{
                Context: ctx,
            }

            stats, err := client.Statistic.GetStatistic(params.Context, params)
            if err != nil {
                return fmt.Errorf("failed to retrieve statistics: %w", err)
            }
			
            statisticview.PrintStatistics(stats.Payload)

            return nil
        },
    }
}
 