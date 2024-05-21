package login

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
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
				Password(true).
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
				Validate(func(str string) error {
					if str == "" {
						return errors.New("credential name cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
	TrimSpaceOfLoginInputs(loginView)
}

func TrimSpaceOfLoginInputs(loginView *LoginView) {
	loginView.Server = strings.TrimSpace(loginView.Server)
	loginView.Username = strings.TrimSpace(loginView.Username)
	loginView.Password = strings.TrimSpace(loginView.Password)
	loginView.Name = strings.TrimSpace(loginView.Name)
	fmt.Println(loginView.Server, loginView.Username, loginView.Password, loginView.Name)
}
