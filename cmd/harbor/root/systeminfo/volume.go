package systeminfo

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/systeminfo"
	"github.com/goharbor/harbor-cli/pkg/views/systeminfo/volume"
)

func GetVolumesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "volumes",
		Short: "Get system volume information",
		Long:  `Retrieve information about the volumes in the Harbor system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := utils.GetClient()
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			
			params := systeminfo.NewGetVolumesParams()
			resp, err := client.Systeminfo.GetVolumes(context.Background(), params)
			if err != nil {
				return fmt.Errorf("failed to get volume info: %w", err)
			}
			
			volume.PrintVolumeInfo(resp.Payload.Storage)
			return nil
		},
	}
}
