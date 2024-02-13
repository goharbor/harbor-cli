package registry

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/api"
	"github.com/goharbor/harbor-cli/internal/pkg/config"
	"github.com/goharbor/harbor-cli/internal/pkg/constants"
	"github.com/spf13/cobra"
)

// GetRegistryCommand creates a new `harbor get registry` command
func GetRegistryCommand() *cobra.Command {
	var opts config.GetRegistryOptions

	cmd := &cobra.Command{
		Use:   "registry [ID]",
		Short: "get registry by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Invalid argument: %s. Expected an integer.\n", args[0])
				return err
			}
			opts.Id = int64(id)

			credentialName, err := cmd.Flags().GetString(constants.CredentialNameOption)
			if err != nil {
				return err
			}
			return api.RunGetRegistry(opts, credentialName, config.OutputType, config.WideOutput)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&config.OutputType, "output", "o", "", "Output type [json/yaml]")
	flags.BoolVarP(&config.WideOutput, "wide", "", false, "Wide output result [true/false]")
	return cmd
}
