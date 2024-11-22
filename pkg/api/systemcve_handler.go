package api

import (
	"strings"
	"time"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/system_cve_allowlist"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/systemcve/update"
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
