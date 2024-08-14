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
     get "github.com/goharbor/harbor-cli/pkg/views/usergroup/get"
)

func UserGroupGetCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "get [groupID]",
        Short: "get user group details",
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

            response, err := api.GetUserGroup(groupId)
            if err != nil {
                if strings.Contains(err.Error(), "404") {
                    log.Errorf("user group not found with id %d", groupId)
                } else {
                    log.Errorf("failed to get user group: %v", err)
                }
                return
            }

            get.DisplayUserGroup(response.Payload)
        },
    }

    return cmd
}
