package registry

import (
	"strconv"
	"sync"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// NewDeleteRegistryCommand creates a new `harbor delete registry` command
func DeleteRegistryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [registryId]",
		Short: "delete registry by id",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			errChan := make(chan error, len(args))
			if len(args) > 0 {
				for _, arg := range args {
					registryId, _ := strconv.ParseInt(arg, 10, 64)
					wg.Add(1)
					go func(registryId int64) {
						defer wg.Done()
						if err := api.DeleteRegistry(registryId); err != nil {
							errChan <- err
						}
					}(registryId)
				}
			} else {
				registryId := prompt.GetRegistryNameFromUser()
				wg.Add(1)

				go func(registryId int64) {
					defer wg.Done()
					if err := api.DeleteRegistry(registryId); err != nil {
						errChan <- err
					}
				}(registryId)
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
				log.Errorf("failed to delete registry: %v", finalErr)
			}
		},
	}

	return cmd
}
