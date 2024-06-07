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

	scannerRegReq := models.ScannerRegistrationReq{
    Name:             &opts.Name,
    Description:      opts.Description,
    Auth:             opts.Auth,
    AccessCredential: opts.AccessCredential,
    URL:              strfmt.URI(opts.URL),
    Disabled:         &opts.Disabled,
    SkipCertVerify:   &opts.SkipCertVerify,
    UseInternalAddr:  &opts.UseInternalAddr,
}

pingScannerReq := models.ScannerRegistrationSettings{
    Name:             &opts.Name,
    Auth:             opts.Auth,
    AccessCredential: opts.AccessCredential,
    URL:              strfmt.URI(opts.URL),
}

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
