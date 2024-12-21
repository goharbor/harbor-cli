package config

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func ListConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List config items",
		Example: `  harbor config list`,
		Long:    `Get information of all CLI config items`,
		Args:    cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := utils.GetCurrentHarborConfig()
			if err != nil {
				logrus.Errorf("Failed to get config: %v", err)
				return
			}

			// Get the output format
			formatFlag := viper.GetString("output-format")
			if formatFlag != "" {
				// Use utils.PrintFormat if available
				err = utils.PrintFormat(config, formatFlag)
				if err != nil {
					logrus.Errorf("Failed to print config: %v", err)
				}
			} else {
				// Default to YAML format
				data, err := yaml.Marshal(config)
				if err != nil {
					logrus.Errorf("Failed to marshal config to YAML: %v", err)
					return
				}
				fmt.Println(string(data))
			}
		},
	}

	return cmd
}
