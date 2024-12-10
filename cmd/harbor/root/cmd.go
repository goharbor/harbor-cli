package root

import (
	"fmt"

	"github.com/goharbor/harbor-cli/cmd/harbor/root/artifact"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/labels"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/registry"
	repositry "github.com/goharbor/harbor-cli/cmd/harbor/root/repository"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/schedule"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/user"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	output  string
	cfgFile string
	verbose bool
)

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "harbor",
		Short:        "Official Harbor CLI",
		SilenceUsage: true,
		Long:         "Official Harbor CLI",
		Example: `
// Base command:
harbor

// Display help about the command:
harbor help
`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Determine if --config was explicitly set
			userSpecifiedConfig := cmd.Flags().Changed("config")
			// Initialize configuration
			utils.InitConfig(cfgFile, userSpecifiedConfig)

			return nil
		},
	}

	root.PersistentFlags().StringVarP(&output, "output-format", "o", "", "Output format. One of: json|yaml")
	root.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.config/harbor-cli/config.yaml)")
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	err := viper.BindPFlag("output-format", root.PersistentFlags().Lookup("output-format"))
	if err != nil {
		fmt.Println(err.Error())
	}

	err = viper.BindPFlag("config", root.PersistentFlags().Lookup("config"))
	if err != nil {
		fmt.Println(err.Error())
	}

	root.AddCommand(
		versionCommand(),
		LoginCommand(),
		project.Project(),
		registry.Registry(),
		repositry.Repository(),
		user.User(),
		artifact.Artifact(),
		HealthCommand(),
		schedule.Schedule(),
		labels.Labels(),
	)

	return root
}
