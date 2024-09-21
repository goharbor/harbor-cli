package create

import (
	"errors"
	"regexp"

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
					} else if !IsValidPassword(str)  {
					 	return errors.New("password is incorrect")
					} else {
					return nil
					}
				}).
				Description("Password should be 8-128 characters long with at least 1 uppercase, 1 lowercase and 1 number."),
			huh.NewInput().
				Title("Comment").
				Value(&createView.Comment),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}


func IsValidPassword(password string) bool {
	if !regexp.MustCompile(`^.{8,128}$`).MatchString(password) {
		return false
	}else if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	} else if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	} else if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	} else{
		return true
	}
}
