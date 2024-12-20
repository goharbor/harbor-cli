package update

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func UpdateUserProfileView(updateView *models.UserProfile) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Email").
				Value(&updateView.Email).
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
				Title("Realname").
				Value(&updateView.Realname).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("real name cannot be empty")
					}
					if isValid := utils.ValidateFL(str); !isValid {
						return errors.New("please enter correct first and last name format, like `Bob Dylan`")
					}
					return nil
				}),
			huh.NewInput().
				Title("Comment").
				Value(&updateView.Comment),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
