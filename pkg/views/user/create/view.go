package create

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
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
					if strings.TrimSpace(str) == "" {
						return errors.New("user name cannot be empty")
					}
					if isVaild := utils.ValidateName(str); !isVaild {
						return errors.New("username cannot contain special characters")
					}
					return nil
				}),
			huh.NewInput().
				Title("Email").
				Value(&createView.Email).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("email cannot be empty or only spaces")
					}
					if isVaild := utils.ValidateEmail(str); !isVaild {
						return errors.New("please enter correct email format")
					}
					return nil
				}),

			huh.NewInput().
				Title("First and Last Name").
				Value(&createView.Realname).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("real name cannot be empty")
					}
					if isValid := utils.ValidateName(str); !isValid {
						return errors.New("please enter correct first and last name format, like `Bob Dylan`")
					}
					return nil
				}),
			huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&createView.Password).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("password cannot be empty or only spaces")
					}
					if err := utils.ValidatePassword(str); err != nil {
						return err
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
