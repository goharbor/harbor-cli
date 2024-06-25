package root

import (
	"fmt"
	"log"
	"os"

	"github.com/goharbor/harbor-cli/cmd/harbor/root/artifact"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/project"
	"github.com/goharbor/harbor-cli/cmd/harbor/root/registry"
	repositry "github.com/goharbor/harbor-cli/cmd/harbor/root/repository"
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

func initConfig() {
	viper.SetConfigType("yaml")

	// cfgFile = viper.GetStering("config")
	viper.SetConfigFile(cfgFile)

	if cfgFile != utils.DefaultConfigPath {
		viper.SetConfigFile(cfgFile)
	} else {
		stat, err := os.Stat(utils.DefaultConfigPath)
		if !os.IsNotExist(err) && stat.Size() == 0 {
			log.Println("Config file is empty, creating a new one")
		}

		if os.IsNotExist(err) {
			log.Printf("Config file not found at %s, creating a new one", cfgFile)
		}

		if os.IsNotExist(err) || (!os.IsNotExist(err) && stat.Size() == 0) {
			if _, err := os.Stat(utils.HarborFolder); os.IsNotExist(err) {
				// Create the parent directory if it doesn't exist

				fmt.Println("Creating config file", utils.HarborFolder)
				if err := os.MkdirAll(utils.HarborFolder, os.ModePerm); err != nil {
					log.Fatal(err)
				}
			}
			err = utils.CreateConfigFile()

			if err != nil {
				log.Fatal(err)
			}

			err = utils.AddCredentialsToConfigFile(utils.Credential{}, cfgFile)

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("Config file created at %s", cfgFile)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

}

func RootCmd() *cobra.Command {
	utils.SetLocation()

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
		// RunE: func(cmd *cobra.Command, args []string) error {

		// },
	}

	cobra.OnInitialize(initConfig)

	root.PersistentFlags().StringVarP(&output, "output-format", "o", "", "Output format. One of: json|yaml")
	root.PersistentFlags().StringVar(&cfgFile, "config", utils.DefaultConfigPath, "config file (default is $HOME/.harbor/config.yaml)")
	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	err := viper.BindPFlag("output-format", root.PersistentFlags().Lookup("output-format"))
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
	)

	return root
}
