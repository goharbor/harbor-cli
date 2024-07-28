package create

import (
	"errors"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	Username string
	Email    string
	Realname string
	Comment  string
	Password string
}

func CreateUserView(createView *CreateView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("User Name").
				Value(&createView.Username).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("user name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Email").
				Value(&createView.Email).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("email cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Title("First and Last Name").
				Value(&createView.Realname).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("real name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
        EchoMode(huh.EchoModePassword).
				Value(&createView.Password).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("password cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Comment").
				Value(&createView.Comment),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
