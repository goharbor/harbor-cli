package registry

import (
	"github.com/spf13/cobra"
)

func Registry() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "registry",
		Short:   "Manage registries",
		Long:    `Manage registries in Harbor`,
		Example: `  harbor registry list`,
	}

	cmd.PersistentFlags().StringP("credential-name", "c", "", "Credential to use of authorization")
	cmd.AddCommand(
		CreateRegistryCommand(),
		DeleteRegistryCommand(),
		ListRegistryCommand(),
		UpdateRegistryCommand(),
		ViewCommand(),
	)

	return cmd
}
