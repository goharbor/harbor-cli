package configurations

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
)

// cmd/harbor/root/config/get.go (or wherever your config get command is)
func GetConfigCmd() *cobra.Command {
	var category string
	var listCategories bool

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Harbor configurations",
		Long: `Get Harbor system configurations. You can filter by category:
- authentication: User and service authentication settings
- security: Security policies and certificate settings  
- system: General system behavior and storage settings`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if listCategories {
				fmt.Println("Available categories:")
				for _, cat := range utils.GetValidCategories() {
					fmt.Printf("  - %s\n", cat)
				}
				return nil
			}

			if category != "" {
				// Get filtered configurations
				configs, err := api.GetConfigurationsByCategory(category)
				if err != nil {
					return err
				}
				return displayCategoryConfigurations(category, configs)
			}

			// Get all configurations
			response, err := api.GetConfigurations()
			if err != nil {
				return err
			}

			// configs := utils.ConvertToConfigurations(response.Payload)
			if err := utils.AddConfigurationsToConfigFile(response.Payload); err != nil {
				return fmt.Errorf("failed to update config file: %v", err)
			}

			return nil
			// displayAllConfigurations(configs)
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter by category (authentication, security, system)")
	cmd.Flags().BoolVar(&listCategories, "list-categories", false, "List available categories")

	return cmd
}

func displayCategoryConfigurations(category string, configs map[string]interface{}) error {
	fmt.Printf("=== %s Configurations ===\n", strings.Title(category))

	if len(configs) == 0 {
		fmt.Printf("No %s configurations found.\n", category)
		return nil
	}

	for key, value := range configs {
		fmt.Printf("%-30s: %v\n", key, formatConfigValue(value))
	}

	return nil
}

func formatConfigValue(value interface{}) string {
	if ptr, ok := value.(*string); ok && ptr != nil {
		return *ptr
	}
	if ptr, ok := value.(*bool); ok && ptr != nil {
		return fmt.Sprintf("%t", *ptr)
	}
	if ptr, ok := value.(*int64); ok && ptr != nil {
		return fmt.Sprintf("%d", *ptr)
	}
	return fmt.Sprintf("%v", value)
}

func displayAllConfigurations(configs *models.Configurations) error {
	fmt.Println("=== All Harbor Configurations ===")

	configValue := reflect.ValueOf(configs).Elem()
	configType := configValue.Type()

	for i := 0; i < configValue.NumField(); i++ {
		fieldName := configType.Field(i).Name
		fieldValue := configValue.Field(i)

		if fieldValue.IsValid() && !fieldValue.IsNil() {
			fmt.Printf("%-30s: %v\n", fieldName, formatConfigValue(fieldValue.Interface()))
		}
	}

	return nil
}
