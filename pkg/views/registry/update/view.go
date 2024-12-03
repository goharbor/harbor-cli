package update

import (
	"errors"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	log "github.com/sirupsen/logrus"
)

func UpdateRegistryView(updateView *models.Registry) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Provider").
				Value(&updateView.Type).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("provider cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Name").
				Value(&updateView.Name).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&updateView.Description),
			huh.NewInput().
				Title("URL").
				Value(&updateView.URL).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("url cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Access ID").
				Value(&updateView.Credential.AccessKey),
			huh.NewInput().
				Title("Access Secret").
				EchoMode(huh.EchoModePassword).
				Description("Replace the Access Secret to the real one").
				Value(&updateView.Credential.AccessSecret),
			huh.NewConfirm().
				Title("Verify Cert").
				Value(&updateView.Insecure).
				Affirmative("yes").
				Negative("no"),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
