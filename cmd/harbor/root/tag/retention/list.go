package retention

import (
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/retention"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/retention/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListExecutionRetentionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list retention execution of the project",
		Args:  cobra.MaximumNArgs(1),
		Example: `harbor retention list [retentionid]`,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var resp retention.ListRetentionExecutionsOK
			var retentionID int
			var strretenId string
			if len(args) > 0 {
				retentionID,_ = strconv.Atoi(args[0])
				resp, err = api.ListRetention(int32(retentionID))
			} else {
				projectId := fmt.Sprintf("%d",prompt.GetProjectIDFromUser())
				strretenId,err = api.GetRetentionId(projectId)
				if err != nil {
					log.Fatal(err)
				}
				retentionID,_ := strconv.Atoi(strretenId)
				resp, err = api.ListRetention(int32(retentionID))
			}

			if err != nil {
				log.Errorf("failed to list retention execution: %v", err)
			}
			FormatFlag := viper.GetString("output-format")
			if FormatFlag != "" {
				utils.PrintPayloadInJSONFormat(resp)
				return
			}

			list.ListRetentionRules(resp.Payload)

		},
	}

	return cmd
}