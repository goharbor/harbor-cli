package login

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type LoginView struct {
	Server   string
	Username string
	Password string
	Name     string
}

func CreateView(loginView *LoginView) {
	theme := huh.ThemeCharm()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Server").
				Description("Server address eg. demo.goharbor.io").
				Value(&loginView.Server).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("server cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("User Name").
				Value(&loginView.Username).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("username cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&loginView.Password).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("password cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Name of Credential").
				Value(&loginView.Name).
				Description("Name of credential to be stored in the harbor config file.").
				PlaceholderFunc(func() string {
					return fmt.Sprintf("%s@%s", loginView.Username, utils.SanitizeServerAddress(loginView.Server))
				}, &loginView).
				SuggestionsFunc(func() []string {
					return []string{
						fmt.Sprintf("%s@%s", loginView.Username, utils.SanitizeServerAddress(loginView.Server)),
					}
				}, &loginView).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("name cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}

}
