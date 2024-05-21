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
				Password(true).
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
	lengthPattern := `^.{8,128}$`
	uppercasePattern := `[A-Z]`
	lowercasePattern := `[a-z]`
	numberPattern := `[0-9]`

	lengthRegex := regexp.MustCompile(lengthPattern)
	uppercaseRegex := regexp.MustCompile(uppercasePattern)
	lowercaseRegex := regexp.MustCompile(lowercasePattern)
	numberRegex := regexp.MustCompile(numberPattern)

	if !lengthRegex.MatchString(password) {
		return false
	}
	if !uppercaseRegex.MatchString(password) {
		return false
	}
	if !lowercaseRegex.MatchString(password) {
		return false
	}
	if !numberRegex.MatchString(password) {
		return false
	}
	return true
}
