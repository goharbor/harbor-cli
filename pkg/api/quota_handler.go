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
package api

import (
	"fmt"
	"strconv"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/quota"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListQuota(opts ListQuotaFlags) (*quota.ListQuotasOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Quota.ListQuotas(
		ctx,
		&quota.ListQuotasParams{
			Page:        &opts.Page,
			PageSize:    &opts.PageSize,
			Reference:   &opts.Reference,
			ReferenceID: &opts.ReferenceID,
			Sort:        &opts.Sort,
		},
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func GetQuotaByRef(projectID int64) (*models.Quota, error) {
	quotas, err := GetAllQuotas(ListQuota, ListQuotaFlags{PageSize: 0})
	if err != nil {
		return nil, err
	}

	// get the correct quota with the Reference
	for _, quota := range quotas {
		pID, err := getRefProjectID(quota.Ref)
		if err != nil {
			return nil, fmt.Errorf("unable to get projectID from quotaRef: %v", err)
		}
		id, _ := strconv.ParseInt(pID, 10, 64)
		if projectID == id {
			return quota, nil
		}
	}
	return nil, fmt.Errorf("unable to find quota with projectID")
}

func GetQuota(quotaID int64) (*models.Quota, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Quota.GetQuota(ctx, &quota.GetQuotaParams{ID: quotaID})
	if err != nil {
		return nil, err
	}

	return response.Payload, nil
}

func UpdateQuota(quotaID int64, hardlimit *models.QuotaUpdateReq) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Quota.UpdateQuota(
		ctx,
		&quota.UpdateQuotaParams{ID: quotaID, Hard: hardlimit},
	)
	if err != nil {
		return err
	}

	return nil
}

// helper function to fetch quotas all quotas
func GetAllQuotas(listFunc func(ListQuotaFlags) (*quota.ListQuotasOK, error), opts ListQuotaFlags) ([]*models.Quota, error) {
	var allQuotas []*models.Quota
	if opts.PageSize == 0 {
		opts.PageSize = 100
		opts.Page = 1

		for {
			quotas, err := listFunc(opts)
			if err != nil {
				return nil, err
			}

			allQuotas = append(allQuotas, quotas.Payload...)

			if len(quotas.Payload) < int(opts.PageSize) {
				break
			}

			opts.Page++
		}
	} else {
		quotas, err := listFunc(opts)
		if err != nil {
			return nil, err
		}
		allQuotas = quotas.Payload
	}

	return allQuotas, nil
}

// helper Function to get project ref details
func getRefProjectID(ref models.QuotaRefObject) (string, error) {
	if refMap, ok := ref.(map[string]interface{}); ok {
		id, _ := refMap["id"]
		return fmt.Sprintf("%v", id), nil
	}
	return "", fmt.Errorf("Error: Ref is not of expected type")
}
