package usergroup

import (
    "bufio"
    "fmt"
    "os"
    "strings"
 
    search "github.com/goharbor/harbor-cli/pkg/views/usergroup/search"
    "github.com/goharbor/harbor-cli/pkg/api"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
)

func UserGroupsSearchCommand() *cobra.Command {
    var groupName string

    cmd := &cobra.Command{
        Use:   "search [groupName]",
        Short: "search user groups",
        Args:  cobra.MaximumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            if len(args) > 0 {
                groupName = args[0]
            }

            if groupName == "" {
                fmt.Print("Enter group name: ")
                reader := bufio.NewReader(os.Stdin)
                input, _ := reader.ReadString('\n')
                groupName = strings.TrimSpace(input)
            }
            fmt.Print("\033[K") 

            fmt.Printf("Searching for groups with name '%s'...\r", groupName)
            response, err := api.SearchUserGroups(groupName)
            if err != nil {
                log.Errorf("failed to search user groups: %v", err)
                return
            }

           
            if len(response.Payload) == 0 {
                log.Infof("No user groups found with the name %s", groupName)
                return
            }

            search.DisplayUserGroupSearchResults(response)
 
        },
    }

    flags := cmd.Flags()
    flags.StringVarP(&groupName, "name", "n", "", "Group name to search")

    return cmd
}
