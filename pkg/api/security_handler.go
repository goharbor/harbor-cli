package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/securityhub"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func GetSecuritySummary() error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	params := &securityhub.GetSecuritySummaryParams{}
	response, err := client.Securityhub.GetSecuritySummary(ctx, params)
	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response.Payload)
	return nil
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
