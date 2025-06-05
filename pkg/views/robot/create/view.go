package create

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	Description string             `json:"description,omitempty"`
	Disable     bool               `json:"disable,omitempty"`
	Duration    int64              `json:"duration,omitempty"`
	Level       string             `json:"level,omitempty"`
	Name        string             `json:"name,omitempty"`
	Permissions []*RobotPermission `json:"permissions"`
	Secret      string             `json:"secret,omitempty"`
	ProjectName string
}

type RobotPermission struct {
	Access    []*models.Access `json:"access"`
	Kind      string           `json:"kind,omitempty"`
	Namespace string           `json:"namespace,omitempty"`
}

type Access struct {
	Action   string `json:"action,omitempty"`
	Effect   string `json:"effect,omitempty"`
	Resource string `json:"resource,omitempty"`
}

func CreateRobotView(createView *CreateView) {
	var duration string
	duration = strconv.FormatInt(createView.Duration, 10)

	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Name").
				Value(&createView.Name).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("Name cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Description").
				Value(&createView.Description),
			huh.NewInput().
				Title("Expiration").
				Value(&duration).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("Expiration cannot be empty")
					}
					dur, err := strconv.ParseInt(str, 10, 64)
					if err != nil {
						return errors.New("invalid expiration time: Enter expiration time in days")
					}
					createView.Duration = dur
					return nil
				}),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func CreateRobotSecretView(name string, secret string) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Robot Name").
				Value(&name),
			huh.NewInput().
				Title("Secret").
				Value(&secret),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}
}
