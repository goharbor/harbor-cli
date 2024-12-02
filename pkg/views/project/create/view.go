package create

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/registry"
	"github.com/goharbor/harbor-cli/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CreateView struct {
	ProjectName  string
	Public       bool
	RegistryID   string
	StorageLimit string
	ProxyCache   bool
}

func getRegistryList() (*registry.ListRegistriesOK, error) {
	credentialName := viper.GetString("current-credential-name")
	client, err := utils.GetClientByCredentialName(credentialName)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	var response *registry.ListRegistriesOK
	response, err = client.Registry.ListRegistries(ctx, &registry.ListRegistriesParams{})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func CreateProjectView(createView *CreateView) {
	theme := huh.ThemeCharm()
	// I want it to be a map of registry ID to registry name
	registries, _ := getRegistryList()

	registryOptions := map[string]string{}
	for _, registry := range registries.Payload {
		regiId := fmt.Sprintf("%d", registry.ID)
		registryOptions[regiId] = fmt.Sprintf("%s (%s)", registry.Name, registry.URL)
	}

	var registrySelectOptions []huh.Option[string]
	for id, name := range registryOptions {
		registrySelectOptions = append(registrySelectOptions, huh.NewOption(name, id))
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Project Name").
				Value(&createView.ProjectName).
				Validate(func(str string) error {
					if strings.TrimSpace(str) == "" {
						return errors.New("project name cannot be empty or only spaces")
					}
					if isValid := utils.ValidateProjectName(str); !isValid {
						return errors.New("please enter correct project name format")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Public").
				Value(&createView.Public).
				Affirmative("yes").
				Negative("no"),
			huh.NewInput().
				Title("Storage Limit").
				Value(&createView.StorageLimit).
				Validate(func(str string) error {
					// Assuming StorageLimit is an int64
					if strings.TrimSpace(str) == "" {
						return errors.New("storage limit cannot be empty or only spaces")
					}
					if err := utils.ValidateStorageLimit(str); err != nil {
						return err
					}
					return nil
				}),

			huh.NewConfirm().
				Title("Proxy Cache").
				Value(&createView.ProxyCache).
				Affirmative("yes").
				Negative("no"),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Validate(func(str string) error {
					if createView.ProxyCache && str == "" {
						return errors.New("registry ID cannot be empty")
					}
					return nil
				}).
				Description("Select a registry to reference when creating the proxy cache project").
				Title("Registry ID").
				Value(&createView.RegistryID).
				Options(registrySelectOptions...),
		).WithHideFunc(func() bool {
			return !createView.ProxyCache || len(registryOptions) == 0
		}),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
