package usergroup

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
 
    "github.com/goharbor/harbor-cli/pkg/api"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
)

func UserGroupDeleteCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "delete [groupID]",
        Short: "delete user group",
        Args:  cobra.MaximumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            var groupId int64
            var err error

            if len(args) == 0 {
                fmt.Print("Enter group ID: ")
                reader := bufio.NewReader(os.Stdin)
                input, _ := reader.ReadString('\n')
                input = strings.TrimSpace(input)
                groupId, err = strconv.ParseInt(input, 10, 64)
                if err != nil {
                    log.Errorf("invalid group ID: %v", err)
                    return
                }
            } else {
                groupId, err = strconv.ParseInt(args[0], 10, 64)
                if err != nil {
                    log.Errorf("invalid group ID: %v", err)
                    return
                }
            }
            response, err := api.ListUserGroups()
            if err != nil {
                log.Errorf("failed to list user groups: %v", err)
                return
            }

            groupExists := false
            for _, group := range response.Payload {
                if group.ID == groupId {
                    groupExists = true
                    break
                }
            }

            if !groupExists {
                log.Errorf("group ID %d not found", groupId)
                return
            }

            err = api.DeleteUserGroup(groupId)
            if err != nil {
                log.Errorf("failed to delete user group: %v", err)
            }
            fmt.Print("\033[K") 

        },
    }

    return cmd
}
