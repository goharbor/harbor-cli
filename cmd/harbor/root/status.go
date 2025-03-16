package root

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/systeminfo"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/user"
	"github.com/goharbor/harbor-cli/cmd/harbor/internal/version"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func StatusCommand() *cobra.Command {
	var longOutput bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the current credential status",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentCredential := viper.GetString("current-credential-name")
			if currentCredential == "" {
				return fmt.Errorf("no active credentials found")
			}

			var registryAddress string
			creds := viper.Get("credentials").([]interface{})
			for _, cred := range creds {
				c := cred.(map[string]interface{})
				if c["name"] == currentCredential {
					registryAddress = c["serveraddress"].(string)
					break
				}
			}

			if registryAddress == "" {
				return fmt.Errorf("registry address not found for current credential: %s", currentCredential)
			}

			ctx, client, err := utils.ContextWithClient()
			if err != nil {
				return fmt.Errorf("failed to create Harbor client: %v", err)
			}

			userInfo, err := client.User.GetCurrentUserInfo(ctx, &user.GetCurrentUserInfoParams{})
			if err != nil {
				return fmt.Errorf("failed to get current user info: %v", err)
			}

			if longOutput {
				// Detailed Output
				isSysAdmin := userInfo.Payload.SysadminFlag
				isAuthAdmin := userInfo.Payload.AdminRoleInAuth

				sysInfo, err := client.Systeminfo.GetSystemInfo(ctx, &systeminfo.GetSystemInfoParams{})
				if err != nil {
					return fmt.Errorf("failed to get system info: %v", err)
				}
				apiVersion := sysInfo.Payload.HarborVersion

				fmt.Println("\nHarbor CLI Status:")
				fmt.Println("==================")
				fmt.Printf("Logged in as: %s\n", userInfo.Payload.Username)
				fmt.Printf("Registry: %s\n", registryAddress)
				fmt.Printf("API Version: %s\n", *apiVersion)
				fmt.Printf("Connected As: %s (%s)\n", userInfo.Payload.Username, roleString(isSysAdmin, isAuthAdmin))

				// Previously logged-in registries
				fmt.Println("\nPreviously Logged in to the following registries:")
				previousRegistriesMap := make(map[string]struct{})
				for _, cred := range creds {
					c := cred.(map[string]interface{})
					if registry, ok := c["serveraddress"].(string); ok {
						previousRegistriesMap[registry] = struct{}{}
					}
				}
				for registry := range previousRegistriesMap {
					fmt.Printf("- %s\n", registry)
				}

				// Add API version and CLI version info
				fmt.Printf("\nCLI Version: %s\n", version.Version)
				fmt.Printf("OS: %s\n", version.System)
			} else {
				// Short Output
				fmt.Printf("CLI currently logged in as %s to registry %s\n", userInfo.Payload.Username, registryAddress)
			}

			return nil
		},
	}

	// Add the --long flag
	cmd.Flags().BoolVarP(&longOutput, "long", "l", false, "Show detailed credential info")

	return cmd
}

// roleString returns the role of the user based on admin flags
func roleString(isSysAdmin, isAuthAdmin bool) string {
	if isSysAdmin {
		return "Harbor Admin"
	}
	return "User"
}
