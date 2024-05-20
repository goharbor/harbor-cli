package create

import (
	"errors"
	"strings"
	"unicode"

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
				Password(true).
				Value(&createView.Password).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("password cannot be empty")
					} else if len(str) < 8  {
					 	return errors.New("password should atleast 8 character long")
					} else if !containsUpperCase(str){
						return errors.New("password should atleast need one uppercase")
					} else if !containsLowerCase(str){
						return errors.New("password should atleast need one lowercase")
					} else if !containsDigits(str){
						return errors.New("password should atleast need one number")
					} else {
					return nil
					}
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

func containsUpperCase(s string) bool {
    return strings.IndexFunc(s, unicode.IsUpper) >= 0
}

func containsLowerCase(s string) bool {
    return strings.IndexFunc(s, unicode.IsLower) >= 0
}

func containsDigits(s string) bool {
    return strings.IndexFunc(s, unicode.IsDigit) >= 0
}
