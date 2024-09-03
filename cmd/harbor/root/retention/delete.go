package retention

import (
	"fmt"
	"strconv"

	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func DeleteRetentionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete retention rule",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var retentionId int
			var strretenId string
			if len(args) > 0 {
				retentionId,_ = strconv.Atoi(args[0])
				err = api.DeleteRetention(int64(retentionId))
			} else {
				projectId := fmt.Sprintf("%d",prompt.GetProjectIDFromUser())
				strretenId,err = api.GetRetentionId(projectId)
				if err != nil {
					log.Fatal(err)
				}
				retentionId,_ = strconv.Atoi(strretenId)
				err = api.DeleteRetention(int64(retentionId))
			}
			if err != nil {
				log.Errorf("failed to delete retention rule: %v", err)
			}
		},
	}

	return cmd
}

