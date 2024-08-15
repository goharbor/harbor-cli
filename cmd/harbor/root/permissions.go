package root

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/permissions" 
)

func permissionsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "permissions",
		Short: "Manage Harbor permissions",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, client, err := utils.ContextWithClient()
			if err != nil {
				fmt.Printf("Error initializing client: %v\n", err)
				os.Exit(1)
			}
	
			permissionsHandler := api.NewPermissionsHandler(client.Permissions)
	
			perms, err := permissionsHandler.GetPermissions(ctx)
			if err != nil {
				fmt.Printf("Failed to get permissions: %v\n", err)
				os.Exit(1)
			}
			permissions.PrintPermissions(perms)
		},
	}
}