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
	"strings"
	"time"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/system_cve_allowlist"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/cveallowlist/update"
	log "github.com/sirupsen/logrus"
)

func ListSystemCve() (system_cve_allowlist.GetSystemCVEAllowlistOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return system_cve_allowlist.GetSystemCVEAllowlistOK{}, err
	}

	response, err := client.SystemCVEAllowlist.GetSystemCVEAllowlist(ctx, &system_cve_allowlist.GetSystemCVEAllowlistParams{})
	if err != nil {
		return system_cve_allowlist.GetSystemCVEAllowlistOK{}, err
	}

	return *response, nil
}

func UpdateSystemCve(opts update.UpdateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}
	var unixTimestamp int64
	if opts.IsExpire {
		expiresAt, err := time.Parse("2006/01/02", opts.ExpireDate)
		if err != nil {
			return err
		}
		unixTimestamp = expiresAt.Unix()
	} else {
		unixTimestamp = 0
	}

	var items []*models.CVEAllowlistItem
	cveIds := strings.Split(opts.CveId, ",")
	for _, id := range cveIds {
		id = strings.TrimSpace(id)
		items = append(items, &models.CVEAllowlistItem{CVEID: id})
	}
	response, err := client.SystemCVEAllowlist.PutSystemCVEAllowlist(ctx, &system_cve_allowlist.PutSystemCVEAllowlistParams{Allowlist: &models.CVEAllowlist{Items: items, ExpiresAt: &unixTimestamp}})
	if err != nil {
		return err
	}

	if response != nil {
		log.Info("cveallowlist added successfully")
	}
	return nil
}
