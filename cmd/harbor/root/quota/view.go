package quota

import (
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/quota"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/views/quota/list"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// View a specified quota
func ViewQuotaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [QuotaID]",
		Short: "get quota by quota ID",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var quota *quota.GetQuotaOK

			if len(args) > 0 {
				quotaID, _ := strconv.ParseInt(args[0], 10, 64)
				quota, err = api.GetQuota(int64(quotaID))
				if err != nil {
					log.Errorf("failed to get Quota: %v", err)
				}
			} else {
				quotaID := prompt.GetQuotaIDFromUser()
				quota, err = api.GetQuota(quotaID)
			}

			if err != nil {
				log.Errorf("failed to get project: %v", err)
			}

			quotas := []*models.Quota{quota.Payload}
			list.ListQuotas(quotas)
		},
	}

	return cmd
}
