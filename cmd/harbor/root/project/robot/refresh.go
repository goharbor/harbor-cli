package robot

import (
	"fmt"
	"os"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/robot/create"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// handle robot view with interactive like in list command.
func RefreshSecretCommand() *cobra.Command {
	var (
		robotID     int64
		secret      string
		secretStdin bool
	)
	cmd := &cobra.Command{
		Use:   "refresh [robotID]",
		Short: "refresh robot secret by id",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if len(args) == 1 {
				robotID, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					log.Errorf("failed to parse robot ID: %v", err)
				}
			} else {
				projectID := prompt.GetProjectIDFromUser()
				robotID = prompt.GetRobotIDFromUser(projectID)
			}

			if secret != "" {
				utils.ValidatePassword(secret)
			}
			if secretStdin {
				secret = getSecret()
			}

			response, err := api.RefreshSecret(secret, robotID)
			if err != nil {
				log.Errorf("failed to refresh robot secret.")
				os.Exit(1)
			}

			log.Info("Secret updated successfully.")

			secret = response.Payload.Secret
			create.CreateRobotSecretView("", secret)

			err = clipboard.WriteAll(response.Payload.Secret)
			fmt.Println("secret copied to clipboard.")
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&secret, "secret", "", "", "secret")
	flags.BoolVarP(&secretStdin, "secret-stdin", "", false, "Take the robot secret from stdin")

	return cmd
}

// getSecret from commandline
func getSecret() string {
	secret, err := utils.GetSecretStdin("Enter your secret: ")
	if err != nil {
		log.Errorf("Error reading secret: %v\n", err)
		os.Exit(1)
	}

	if err := utils.ValidatePassword(secret); err != nil {
		log.Errorf("Invalid secret: %v\n", err)
		os.Exit(1)
	}
	return secret
}
