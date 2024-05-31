package api

import (
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

func GetQuota(quotaID int64) (*quota.GetQuotaOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Quota.GetQuota(ctx, &quota.GetQuotaParams{ID: quotaID})
	if err != nil {
		return nil, err
	}

	return response, nil
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
