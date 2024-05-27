package user

import (
	"strconv"
	"sync"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func UserDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete user by UserID",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			errChan := make(chan error, len(args)) // Channel to collect error

			if len(args) > 0 {
				for _, arg := range args {
					userId, _ := strconv.ParseInt(arg, 10, 64)
					wg.Add(1)
					go func(userId int64) {
						defer wg.Done()
						if err := api.DeleteUser(userId); err != nil {
							errChan <- err
						}
					}(userId)
				}
			} else {
				userId := prompt.GetUserIdFromUser()
				wg.Add(1)
				go func(userId int64) {
					defer wg.Done()
					if err := api.DeleteUser(userId); err != nil {
						errChan <- err
					}
				}(userId)
			}

			// Wait for all goroutines to finish
			go func() {
				wg.Wait()
				close(errChan)
			}()

			// Collect and handle errors
			var finalErr error
			for err := range errChan {
				if finalErr == nil {
					finalErr = err
				} else {
					log.Errorf("Error: %v", err)
				}
			}
			if finalErr != nil {
				log.Errorf("failed to delete user: %v", finalErr)
			}
		},
	}

	return cmd
}
