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
	"fmt"
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
		storage string
	)

	var opts api.ListQuotaFlags
	cmd := &cobra.Command{
		Use:   "update [QuotaID]",
		Short: "update quotas for projects",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var storageValue int64

			// get quota id with quota
			quota, err := GetQuotaFromUser(args, opts)
			if err != nil {
				log.Errorf("error: %v", err)
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
				storage = update.UpdateQuotaView(quota)
				storageValue, err = utils.StorageStringToBytes(storage)
				if err != nil {
					log.Errorf("failed to parse storage: %v", err)
					os.Exit(1)
				}
			}

			hardlimit := &models.QuotaUpdateReq{
				Hard: models.ResourceList{"storage": storageValue},
			}

			err = api.UpdateQuota(quota.ID, hardlimit)
			if err != nil {
				log.Errorf("failed to update quota: %v", err)
				os.Exit(1)
			}

			log.Infof("quota updated successfully!")
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&storage, "storage", "", "", "Enter storage size (e.g., 50GiB, 20MiB, 4TiB)")
	flags.StringVarP(&opts.Reference, "project-name", "", "", "Get quota by project-name")
	flags.StringVarP(&opts.ReferenceID, "project-id", "", "", "Get quota by project ID")

	return cmd
}

func GetQuotaFromUser(args []string, opts api.ListQuotaFlags) (*models.Quota, error) {
	var err error
	var quota *models.Quota

	if len(args) > 0 {
		quotaID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			err := fmt.Errorf("failed to parse quotaID: %v", err)
			return nil, err
		}
		quota, err = api.GetQuota(int64(quotaID))
		if err != nil {
			err := fmt.Errorf("failed to get Quota: %v", err)
			return nil, err
		}
	} else if opts.Reference != "" {
		project, err := api.GetProject(opts.Reference, false)
		if err != nil {
			err := fmt.Errorf("failed to get project: %v", err)
			return nil, err
		}
		projectID := project.Payload.ProjectID
		quota, err = api.GetQuotaByRef(int64(projectID))
		if err != nil {
			err := fmt.Errorf("failed to get quota: %v", err)
			return nil, err
		}
	} else if opts.ReferenceID != "" {
		projectID, err := strconv.ParseInt(opts.ReferenceID, 10, 64)
		if err != nil {
			err := fmt.Errorf("invalid projectID: %v", err)
			return nil, err
		}
		quota, err = api.GetQuotaByRef(projectID)
		if err != nil {
			err := fmt.Errorf("failed to get quota: %v", err)
			return nil, err
		}
	} else {
		quotaID := prompt.GetQuotaIDFromUser()
		if quotaID == 0 {
			err := fmt.Errorf("failed to get quotaID from user")
			return nil, err
		}
		quota, err = api.GetQuota(quotaID)
		if err != nil {
			err := fmt.Errorf("failed to get quota: %v", err)
			return nil, err
		}
	}

	return quota, nil
}
