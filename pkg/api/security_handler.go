package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/securityhub"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func GetSecuritySummary() (*securityhub.GetSecuritySummaryOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Securityhub.GetSecuritySummary(ctx,&securityhub.GetSecuritySummaryParams{})
	if err != nil {
		return nil,err
	}

	return response,nil
}

func ListVulnerabilities(query string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	params := &securityhub.ListVulnerabilitiesParams{
		Q: &query,
	}

	response, err := client.Securityhub.ListVulnerabilities(ctx, params)
	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.Payload)
	return nil
}
