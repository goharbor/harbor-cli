package api

import (
	"fmt"
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

	// The input for auth is not clearly stated in the API docs
	// The source code states it needs an empty string https://github.com/goharbor/harbor/blob/115827cac7eb5753160b63339f48108937ed673e/src/pkg/scan/dao/scanner/model.go#L50C2-L50C19
	if opts.Auth == "None" {
		opts.Auth = ""
	}

	url := strfmt.URI(opts.URL)
	scannerRegReq := models.ScannerRegistrationReq{
		Name:             &opts.Name,
		Description:      opts.Description,
		Auth:             opts.Auth,
		AccessCredential: opts.AccessCredential,
		URL:              &url,
		Disabled:         &opts.Disabled,
		SkipCertVerify:   &opts.SkipCertVerify,
		UseInternalAddr:  &opts.UseInternalAddr,
	}

	_, err = client.Scanner.CreateScanner(ctx, &scanner.CreateScannerParams{Registration: &scannerRegReq})

	if err != nil {
		return err
	} else {
		fmt.Println("Scanner created successfully.")
	}

	return nil
}

func ListScanners() (scanner.ListScannersOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return scanner.ListScannersOK{}, err
	}

	response, err := client.Scanner.ListScanners(ctx, &scanner.ListScannersParams{})

	if err != nil {
		return scanner.ListScannersOK{}, err
	}

	return *response, nil
}

func GetScanner(registrationID string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.Scanner.GetScanner(ctx, &scanner.GetScannerParams{RegistrationID: registrationID})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}

func GetScannerMetadata(registrationID string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	response, err := client.Scanner.GetScannerMetadata(ctx, &scanner.GetScannerMetadataParams{RegistrationID: registrationID})

	if err != nil {
		return err
	}

	utils.PrintPayloadInJSONFormat(response)
	return nil
}

func SetDefaultScanner(registrationID string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Scanner.SetScannerAsDefault(ctx, &scanner.SetScannerAsDefaultParams{RegistrationID: registrationID, Payload: &models.IsDefault{IsDefault: true}})

	if err != nil {
		return err
	}

	return nil
}

func DeleteScanner(registrationID string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Scanner.DeleteScanner(ctx, &scanner.DeleteScannerParams{RegistrationID: registrationID})

	if err != nil {
		return err
	}

	return nil
}

func UpdateScanner(registrationID string, opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	url := strfmt.URI(opts.URL)
	scannerRegReq := models.ScannerRegistrationReq{
		Name:             &opts.Name,
		Description:      opts.Description,
		Auth:             opts.Auth,
		AccessCredential: opts.AccessCredential,
		URL:              &url,
		Disabled:         &opts.Disabled,
		SkipCertVerify:   &opts.SkipCertVerify,
		UseInternalAddr:  &opts.UseInternalAddr,
	}

	_, err = client.Scanner.UpdateScanner(ctx, &scanner.UpdateScannerParams{RegistrationID: registrationID, Registration: &scannerRegReq})

	if err != nil {
		return err
	} else {
		fmt.Println("Scanner updated successfully.")
	}

	return nil

}

func PingScanner(opts create.CreateView) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	url := strfmt.URI(opts.URL)
	scannerPingReq := models.ScannerRegistrationSettings{
		Name:             &opts.Name,
		Auth:             opts.Auth,
		AccessCredential: opts.AccessCredential,
		URL:              &url,
	}

	_, err = client.Scanner.PingScanner(ctx, &scanner.PingScannerParams{Settings: &scannerPingReq})

	if err != nil {
		return err
	} else {
		fmt.Println("Scanner pinged successfully.")
	}

	return nil
}
