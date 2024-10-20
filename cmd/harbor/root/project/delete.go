package project

import (
	"sync"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DeleteProjectCommand creates a new `harbor delete project` command
func DeleteProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [NAME|ID]",
		Short:   "delete project by name or id",
		Example: `  harbor project delete [NAME|ID]`,
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			errChan := make(
				chan error,
				len(args),
			) // Channel to collect errors

			if len(args) > 0 {
				for _, arg := range args {
					wg.Add(1)
					go func(projectName string) {
						defer wg.Done()
						if err := api.DeleteProject(projectName); err != nil {
							errChan <- err
						}
					}(arg)
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
				log.Errorf("failed to delete some projects: %v", finalErr)
			}
		},
	}

	return cmd
}
