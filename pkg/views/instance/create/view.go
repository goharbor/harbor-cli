package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct{
	Vendor	string
	Name	string
	Description	string
	Endpoint	string
	AuthMode string
	AuthInfo map[string]string
	Enabled bool
	Insecure bool
}


func CreateInstanceView(createView *CreateView) {
	cv := CreateView{
        AuthInfo: map[string]string{
            "username": "",
            "password": "",
			"token":	"",
        },
    }
	username := cv.AuthInfo["username"]
	password := cv.AuthInfo["password"]
	token := cv.AuthInfo["token"]
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Provider").
				Options(
					huh.NewOption("Dragonfly", "dragonfly"),
					huh.NewOption("Kraken", "kraken"),
				).
				Value(&createView.Vendor),
			huh.NewInput().
				Title("Name").
				Value(&createView.Name).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("instance name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description),
			huh.NewInput().
				Title("Endpoint").
				Value(&createView.Endpoint).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("endpoint cannot be empty")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Enable").
				Value(&createView.Enabled).
				Affirmative("yes").
				Negative("no"),
			huh.NewConfirm().
				Title("Verify Cert").
				Value(&createView.Insecure).
				Affirmative("yes").
				Negative("no"),
			huh.NewSelect[string]().
				Title("Auth Mode").
				Options(
					huh.NewOption("None", "NONE"),
					huh.NewOption("Basic", "BASIC"),
					huh.NewOption("OAuth", "OAUTH"),
				).
				Value(&createView.AuthMode),
			),
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Value(&username).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("username cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				Value(&password).
				Password(true).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("password cannot be empty")
					}
					return nil
				}),
				).WithHideFunc(func() bool {
					return createView.AuthMode == "NONE" || createView.AuthMode == "OAUTH"
				}),
		huh.NewGroup(
			huh.NewInput().
				Title("Token").
				Value(&token).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("token cannot be empty")
					}
					return nil
				}),
				).WithHideFunc(func() bool {
					return createView.AuthMode == "NONE" || createView.AuthMode == "BASIC"
				}),
		).WithTheme(theme).Run()
		
	if err != nil {
		log.Fatal(err)
	}
}