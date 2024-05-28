package member

import (
	"context"
	"strconv"
	"sync"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/member"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Deletes the member of the given project and Member
func DeleteMemberCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [projectName or ID] [memberID]",
		Short: "delete member by id",
		Args:  cobra.MinimumNArgs(0),
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

			if args[1] == "%" {
				runDeleteAllMember(args[0])
			} else if len(args) > 1 {
				for _, mid := range memberID {
					wg.Add(1)
					go func(member int64) {
						defer wg.Done()
						err := runDeleteMember(args[0], member)
						if err != nil {
							errChan <- err
						}
					}(mid)
				}
			} else {
				projectName := utils.GetProjectNameFromUser()
				memID := utils.GetMemberIDFromUser(projectName)
				wg.Add(1)
				go func(member int64) {
					defer wg.Done()
					err := runDeleteMember(projectName, memID)
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

func runDeleteAllMember(projectName string) {
	var wg sync.WaitGroup
	errChan := make(chan error, 0)
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	response, err := client.Member.ListProjectMembers(
		ctx,
		&member.ListProjectMembersParams{ProjectNameOrID: projectName},
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, member := range response.Payload {
		wg.Add(1)
		go func(memberID int64) {
			defer wg.Done()
			err := runDeleteMember(projectName, memberID)
			if err != nil {
				errChan <- err
			}
		}(member.ID) // Pass member.ID to the goroutine
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Handle errors after all deletions are done
	for err := range errChan {
		if err != nil {
			log.Errorln("Error:", err)
		}
	}
}

func runDeleteMember(projectName string, memberID int64) error {
	credentialName := viper.GetString("current-credential-name")
	client := utils.GetClientByCredentialName(credentialName)
	ctx := context.Background()
	_, err := client.Member.DeleteProjectMember(
		ctx,
		&member.DeleteProjectMemberParams{ProjectNameOrID: projectName, Mid: memberID},
	)
	if err != nil {
		return err
	}

	log.Info("Member deleted successfully")
	return nil
}
