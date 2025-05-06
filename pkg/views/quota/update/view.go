// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package update

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/quota/list"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	StorageUnit string
	Value       int64
}

func UpdateQuotaView(quta *models.Quota) string {
	var (
		value      string
		createView CreateView
	)

	rawStorage := list.BytesToStorageString(quta.Hard["storage"])
	fmt.Printf("current storage: %v", rawStorage)
	storagearr := strings.Split(rawStorage, " ")
	storageUnits := []string{"MiB", "GiB", "TiB"}
	storageUnit := storagearr[1]

	// check if storageUnit is valid
	if slices.Contains(storageUnits, storageUnit) {
		createView.StorageUnit = storageUnit
		value = storagearr[0]
	}

	// Initialize a slice to hold select options
	var storageUnitSelectOptions []huh.Option[string]

	// Iterate over storageUnits to populate storageUnitSelectOptions
	for _, option := range storageUnits {
		storageUnitSelectOptions = append(
			storageUnitSelectOptions,
			huh.NewOption(option, option),
		)
	}

	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a Storage Unit").
				Value(&createView.StorageUnit).
				Options(storageUnitSelectOptions...).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("Storage Type cannot be empty.")
					}
					return nil
				}),

			huh.NewInput().
				Title("Quota Limit").
				Value(&value).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("Quota Limits cannot be empty")
					}
					intval, err := strconv.ParseInt(str, 10, 64)
					if err != nil {
						return errors.New("Quota limit must be a valid integer")
					}
					createView.Value = intval
					return nil
				}),
		),
	).WithTheme(theme).Run()
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%v%v", createView.Value, createView.StorageUnit)
}
