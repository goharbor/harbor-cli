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
		Short: "Ping the server",
		Long:  "Ping the server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Pinging...")

			err := api.Ping()
			if err != nil {
				log.Errorf("Failed to ping the server: %v", err)
				return
			}

			fmt.Println("Successfully pinged the server!")
		},
	}

	return cmd
}
