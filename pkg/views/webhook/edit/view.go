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
package edit

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type EditView struct {
	WebhookId               int64
	ProjectName             string
	Name                    string
	Description             string
	NotifyType              string
	PayloadFormat           string
	EventType               []string
	EndpointURL             string
	AuthHeader              string
	VerifyRemoteCertificate bool
	Enabled                 bool
}

func isSelected(selected []string, option string) bool {
	for _, item := range selected {
		if item == option {
			return true
		}
	}
	return false
}

func WebhookEditView(editView *EditView) {
	theme := huh.ThemeCharm()
	var verifyCert string
	var enable string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&editView.Name).
				Validate(utils.EmptyStringValidator("Webhook Name")),

			huh.NewText().
				Title("Description").
				Value(&editView.Description),

			huh.NewSelect[string]().
				Title("Notify Type").
				Options(
					huh.NewOption("http", "http").Selected(editView.NotifyType == "http"),
					huh.NewOption("slack", "slack").Selected(editView.NotifyType == "slack"),
				).
				Value(&editView.NotifyType),

			huh.NewSelect[string]().Title("Webhook Enabled").
				Description("Determine whether the webhook should verify the certificate of a remote url "+
					"Uncheck this box when the remote url uses a self-signed or untrusted certificate.").
				Options(
					huh.NewOption("True", "yes"),
					huh.NewOption("False", "no"),
				).Value(&enable),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}

	if editView.NotifyType == "http" {
		err = huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Payload Format").
					Options(
						huh.NewOption("Default", "Default").Selected(editView.PayloadFormat == "Default"),
						huh.NewOption("CloudEvents", "CloudEvents").Selected(editView.PayloadFormat == "CloudEvents"),
					).
					Value(&editView.PayloadFormat),
			),
		).WithTheme(theme).Run()

		if err != nil {
			log.Fatal(err)
		}
	}

	err = huh.NewForm(
		huh.NewGroup(

			huh.NewInput().Title("Endpoint URL").
				Value(&editView.EndpointURL).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("endpoint URL cannot be empty")
					}
					if err := utils.ValidateURL(str); err != nil {
						return err
					}
					return nil
				}),

			huh.NewInput().
				Title("Auth Header").
				Value(&editView.AuthHeader),

			huh.NewMultiSelect[string]().
				Title("Select Event Types").
				Options(
					huh.NewOption("Artifact deleted", "DELETE_ARTIFACT").
						Selected(isSelected(editView.EventType, "DELETE_ARTIFACT")),
					huh.NewOption("Artifact pulled", "PULL_ARTIFACT").
						Selected(isSelected(editView.EventType, "PULL_ARTIFACT")),
					huh.NewOption("Artifact pushed", "PUSH_ARTIFACT").
						Selected(isSelected(editView.EventType, "PUSH_ARTIFACT")),
					huh.NewOption("Quota exceed", "QUOTA_EXCEED").
						Selected(isSelected(editView.EventType, "QUOTA_EXCEED")),
					huh.NewOption("Quota near threshold", "QUOTA_WARNING").
						Selected(isSelected(editView.EventType, "QUOTA_WARNING")),
					huh.NewOption("Replication status changed", "REPLICATION").
						Selected(isSelected(editView.EventType, "REPLICATION")),
					huh.NewOption("Scanning failed", "SCANNING_FAILED").
						Selected(isSelected(editView.EventType, "SCANNING_FAILED")),
					huh.NewOption("Scanning finished", "SCANNING_COMPLETED").
						Selected(isSelected(editView.EventType, "SCANNING_COMPLETED")),
					huh.NewOption("Scanning stopped", "SCANNING_STOPPED").
						Selected(isSelected(editView.EventType, "SCANNING_STOPPED")),
					huh.NewOption("Tag retention finished", "TAG_RETENTION").
						Selected(isSelected(editView.EventType, "TAG_RETENTION")),
				).
				Value(&editView.EventType).
				Validate(func(args []string) error {
					if len(args) == 0 {
						return errors.New("please select least one of event type(s)")
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

	editView.VerifyRemoteCertificate = (verifyCert == "yes")
	editView.Enabled = (enable == "yes")
	if editView.NotifyType == "slack" {
		editView.PayloadFormat = ""
	}

	if err != nil {
		log.Fatal(err)
	}
}
