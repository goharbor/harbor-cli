package systeminfo

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/systeminfo"
	"github.com/goharbor/harbor-cli/pkg/views/systeminfo/info"
)

func newGetInfoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Get general system information",
		Long:  `Retrieve general information about the Harbor system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := utils.GetClient()
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			
			params := systeminfo.NewGetSystemInfoParams()
			resp, err := client.Systeminfo.GetSystemInfo(context.Background(), params)
			if err != nil {
				return fmt.Errorf("failed to get system info: %w", err)
			}
			
			info.PrintSystemInfo(resp.Payload)   
			return nil
		},
	}
}
