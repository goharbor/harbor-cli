package api

import (
	"github.com/go-openapi/strfmt"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/scanner"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	"github.com/goharbor/harbor-cli/pkg/views/scanner/create"
)

func CreateScanner(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	var scannerRegReq models.ScannerRegistrationReq
	var pingScannerReq models.ScannerRegistrationSettings

	scannerRegReq.Name = &opts.Name
	scannerRegReq.Description = opts.Description
	scannerRegReq.Auth = opts.Auth
	scannerRegReq.AccessCredential = opts.AccessCredential
	url := strfmt.URI(opts.URL)
	scannerRegReq.URL = &url
	scannerRegReq.Disabled = &opts.Disabled
	scannerRegReq.SkipCertVerify = &opts.SkipCertVerify
	scannerRegReq.UseInternalAddr = &opts.UseInternalAddr

	pingScannerReq.Name = &opts.Name
	pingScannerReq.Auth = opts.Auth
	pingScannerReq.AccessCredential = opts.AccessCredential
	pingScannerReq.URL = &url

	_, err = client.Scanner.PingScanner(ctx, &scanner.PingScannerParams{Settings: &pingScannerReq})
	if err != nil {
		return err
	}

	response, err := client.Scanner.CreateScanner(ctx, &scanner.CreateScannerParams{Registration: &scannerRegReq})

	if err != nil {
		return err
	}

	if response != nil {
		return nil
	}
	return nil
}
