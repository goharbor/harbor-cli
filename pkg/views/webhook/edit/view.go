package edit

import (
	"errors"
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

			huh.NewConfirm().Title("Webhook Enabled").
				Description("Determine whether the webhook should verify the certificate of a remote url "+
					"Uncheck this box when the remote url uses a self-signed or untrusted certificate.").
				Affirmative("True").
				Negative("False").
				Value(&editView.Enabled),
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
				Validate(utils.EmptyStringValidator("Endpoint URL")),

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

			huh.NewConfirm().Title("Verify Remote Certificate").
				Description("Determine whether the webhook should verify the certificate of a remote url "+
					"Uncheck this box when the remote url uses a self-signed or untrusted certificate.").
				Affirmative("Yes").
				Negative("No").
				Value(&editView.VerifyRemoteCertificate),
		),
	).WithTheme(theme).Run()

	if editView.NotifyType == "slack" {
		editView.PayloadFormat = ""
	}

	if err != nil {
		log.Fatal(err)
	}
}
