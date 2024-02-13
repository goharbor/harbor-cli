package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "0.1.0"
	GitCommit = ""
)

func VersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of harbor-cli",
		RunE: func(cmd *cobra.Command, args []string) error{
			fmt.Println("Harbor-cli")
			fmt.Println("Version:", Version)
			fmt.Println("Commit:", GitCommit)
			return nil
		},
	}
	return cmd
}