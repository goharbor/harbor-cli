package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/quota"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListQuota(opts ListQuotaFlags) (quota.ListQuotasOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return quota.ListQuotasOK{}, err
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
		return quota.ListQuotasOK{}, err
	}
	return *response, nil
}
