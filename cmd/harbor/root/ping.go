package root

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func PingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Ping the Harbor server",
		Long:  `Check connectivity to the Harbor server. Returns an error if Harbor is unreachable.`,
		Example: `  harbor ping

  # Example usage with verbose output
  harbor ping -v
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := api.Ping()
			if err != nil {
				log.Errorf("Failed to ping Harbor: %v", err)
				return err
			}

			fmt.Println("Pong! Successfully contacted the Harbor server.")
			return nil
		},
	}

	return cmd
}
