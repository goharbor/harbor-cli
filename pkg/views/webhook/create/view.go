package create

import (
	"errors"
	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
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

func WebhookCreateView(createView *CreateView) {
	theme := huh.ThemeCharm()
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
		log.Fatal(err)
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
			log.Fatal(err)
		}
	}

	err = huh.NewForm(
		huh.NewGroup(

			huh.NewInput().Title("Endpoint URL").
				Value(&createView.EndpointURL).
				Validate(utils.EmptyStringValidator("Endpoint URL")),

			huh.NewInput().Title("Auth Header").Value(&createView.AuthHeader),

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
						return errors.New("please select least one of event type(s)")
					}
					return nil
				}),

			huh.NewConfirm().Title("Verify Remote Certificate").
				Description("Determine whether the webhook should verify the certificate of a remote url "+
					"Uncheck this box when the remote url uses a self-signed or untrusted certificate.").
				Affirmative("Yes").
				Negative("No").
				Value(&createView.VerifyRemoteCertificate),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
