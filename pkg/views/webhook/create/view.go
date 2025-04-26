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
package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

type CreateView struct {
	ProjectName             string
	Name                    string
	Description             string
	NotifyType              string
	PayloadFormat           string
	EventType               []string
	EndpointURL             string
	AuthHeader              string
	VerifyRemoteCertificate bool
}

func WebhookCreateView(createView *CreateView) error {
	theme := huh.ThemeCharm()
	var verifyCert string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&createView.Name).
				Validate(utils.EmptyStringValidator("Webhook Name")),

			huh.NewText().
				Title("Description").
				Value(&createView.Description),

			huh.NewSelect[string]().
				Title("Notify Type").
				Options(
					huh.NewOption("http", "http"),
					huh.NewOption("slack", "slack"),
				).
				Value(&createView.NotifyType),
		),
	).WithTheme(theme).Run()

	if err != nil {
		return err
	}

	if createView.NotifyType == "http" {
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Payload Format").
					Options(
						huh.NewOption("Default", "Default"),
						huh.NewOption("CloudEvents", "CloudEvents"),
					).
					Value(&createView.PayloadFormat),
			),
		).WithTheme(theme).Run()

		if err != nil {
			return err
		}
	}

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Endpoint URL").
				Value(&createView.EndpointURL).
				Validate(func(str string) error {
					return utils.ValidateURL(str)
				}),
			huh.NewInput().
				Title("Auth Header").
				Value(&createView.AuthHeader),
			huh.NewMultiSelect[string]().
				Title("Select Event Types").
				Options(
					huh.NewOption("Artifact deleted", "DELETE_ARTIFACT"),
					huh.NewOption("Artifact pulled", "PULL_ARTIFACT"),
					huh.NewOption("Artifact pushed", "PUSH_ARTIFACT"),
					huh.NewOption("Quota exceed", "QUOTA_EXCEED"),
					huh.NewOption("Quota near threshold", "QUOTA_WARNING"),
					huh.NewOption("Replication status changed", "REPLICATION"),
					huh.NewOption("Scanning failed", "SCANNING_FAILED"),
					huh.NewOption("Scanning finished", "SCANNING_COMPLETED"),
					huh.NewOption("Scanning stopped", "SCANNING_STOPPED"),
					huh.NewOption("Tag retention finished", "TAG_RETENTION"),
				).
				Value(&createView.EventType).
				Validate(func(args []string) error {
					if len(args) == 0 {
						return errors.New("please select at least one event type")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("Verify Remote Certificate").
				Description("Determine whether the webhook should verify the certificate of a remote URL.\n"+
					"Uncheck this box when the remote URL uses a self-signed or untrusted certificate.").
				Options(
					huh.NewOption("Yes", "yes"),
					huh.NewOption("No", "no"),
				).
				Value(&verifyCert),
		),
	).WithTheme(theme).Run()

	createView.VerifyRemoteCertificate = (verifyCert == "yes")
	return err
}
