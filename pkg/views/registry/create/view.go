package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	Name        string
	Type        string
	Description string
	URL         string
	Credential  RegistryCredential
	Insecure    bool
}

type RegistryCredential struct {
	AccessKey    string `json:"access_key,omitempty"`
	Type         string `json:"type,omitempty"`
	AccessSecret string `json:"access_secret,omitempty"`
}

func CreateRegistryView(createView *CreateView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Provider").
				Value(&createView.Type).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("provider cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Name").
				Value(&createView.Name).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description),
			huh.NewInput().
				Title("URL").
				Value(&createView.URL).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("url cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Access Key").
				Value(&createView.Credential.AccessKey),
			huh.NewInput().
				Title("Access Secret").
				Value(&createView.Credential.AccessSecret),
			huh.NewConfirm().
				Title("Verify Cert").
				Value(&createView.Insecure).
				Affirmative("yes").
				Negative("no"),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
