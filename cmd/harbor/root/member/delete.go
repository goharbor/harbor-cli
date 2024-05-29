package member

import (
	"strconv"
	"sync"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Deletes the member of the given project and Member
func DeleteMemberCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [projectName or ID] [memberID]",
		Short:   "delete member by id",
		Long:    "delete members in a project by MemberID",
		Example: "  harbor member delete my-project 2",
		Args:    cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var wg sync.WaitGroup
			errChan := make(
				chan error,
				len(args),
			) // Channel to collect errors

			var memberID []int64

			for i, mid := range args {
				if i != 0 {
					mID, _ := strconv.Atoi(mid)
					memberID = append(memberID, int64(mID))
				}
			}

			if len(args) > 1 {
				if args[1] == "%" {
					api.DeleteAllMember(args[0])
				}
				for _, mid := range memberID {
					wg.Add(1)
					go func(member int64) {
						defer wg.Done()
						err := api.DeleteMember(args[0], member)
						if err != nil {
							errChan <- err
						}
					}(mid)
				}
			} else {
				var projectName string
				if len(args) > 0 {
					projectName = args[0]
				} else {
					projectName = prompt.GetProjectNameFromUser()
				}
				memID := prompt.GetMemberIDFromUser(projectName)
				wg.Add(1)
				go func(member int64) {
					defer wg.Done()
					err := api.DeleteMember(projectName, memID)
					if err != nil {
						errChan <- err
					}
				}(memID)
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
				log.Errorf("failed to delete some members: %v", finalErr)
			}
		},
	}

	return cmd
}
