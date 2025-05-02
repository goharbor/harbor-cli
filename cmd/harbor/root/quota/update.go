// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package quota

import (
	"os"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/api"
	"github.com/goharbor/harbor-cli/pkg/prompt"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/quota/update"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type QuotaUpdateReq struct {
	// The new hard limits for the quota
	Hard ResourceList `json:"hard,omitempty"`
}

type ResourceList map[string]int64

// UpdateQuotaCommand updates the quota
func UpdateQuotaCommand() *cobra.Command {
	var (
		quotaID int64
		storage string
	)

	cmd := &cobra.Command{
		Use:   "update [QuotaID]",
		Short: "update quotas for projects",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var storageValue int64

			if len(args) > 0 {
				quotaID, err = strconv.ParseInt(args[0], 10, 64)
			} else {
				quotaID = prompt.GetQuotaIDFromUser()
			}

			if err != nil {
				log.Errorf("failed to parse quotaID: %v", err)
				return
			}

			if storage != "" {
				if storage == "-1" {
					storageValue = -1
				} else {
					storageValue, err = utils.StorageStringToBytes(storage)
					if err != nil {
						log.Errorf("failed to parse storage: %v", err)
						os.Exit(1)
					}
				}
			} else {
				storage = update.UpdateQuotaView()
				storageValue, err = utils.StorageStringToBytes(storage)
				if err != nil {
					log.Errorf("failed to parse storage: %v", err)
					os.Exit(1)
				}
			}

			hardlimit := &models.QuotaUpdateReq{
				Hard: models.ResourceList{"storage": storageValue},
			}

			err = api.UpdateQuota(quotaID, hardlimit)
			if err != nil {
				log.Errorf("failed to update quota: %v", err)
				os.Exit(1)
			}

			log.Infof("quota updated successfully!")
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&storage, "storage", "", "", "Enter storage size (e.g., 50GiB, 200MiB, 4TiB)")

	return cmd
}
