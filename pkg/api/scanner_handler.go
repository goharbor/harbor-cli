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

func GetScanner(registrationID string) (scanner.GetScannerOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return scanner.GetScannerOK{}, err
	}

	response, err := client.Scanner.GetScanner(ctx, &scanner.GetScannerParams{RegistrationID: registrationID})

	if err != nil {
		return scanner.GetScannerOK{}, err
	}

	return *response, nil
}

func GetScannerMetadata(registrationID string) (scanner.GetScannerMetadataOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return scanner.GetScannerMetadataOK{}, err
	}

	response, err := client.Scanner.GetScannerMetadata(ctx, &scanner.GetScannerMetadataParams{RegistrationID: registrationID})

	if err != nil {
		return scanner.GetScannerMetadataOK{}, err
	}

	return *response, nil
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

func UpdateScanner(registrationID string, opts models.ScannerRegistration) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	scannerRegReq := models.ScannerRegistrationReq{
		Name:             &opts.Name,
		Description:      opts.Description,
		Auth:             opts.Auth,
		AccessCredential: opts.AccessCredential,
		URL:              &opts.URL,
		Disabled:         opts.Disabled,
		SkipCertVerify:   opts.SkipCertVerify,
		UseInternalAddr:  opts.UseInternalAddr,
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

func GetScannerByName(name string) (models.ScannerRegistration, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return models.ScannerRegistration{}, err
	}

	response, err := client.Scanner.ListScanners(ctx, &scanner.ListScannersParams{})

	if err != nil {
		return models.ScannerRegistration{}, err
	}

	for _, scanner := range response.GetPayload() {
		if scanner.Name == name {
			return *scanner, nil
		}
	}

	return models.ScannerRegistration{}, fmt.Errorf("scanner with name %s not found", name)
}
