package create

import (
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
)

func CreateTagView(tagName *string) {
	theme := huh.ThemeCharm()

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Tag Name").
				Value(tagName).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("tag name cannot be empty or only spaces")
					}
					if isVaild := utils.VaildateTagName(str); !isVaild {
						return errors.New("please enter the correct tag name format")
					}
					return nil
				}),
		),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
