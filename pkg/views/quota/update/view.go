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
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/views/quota/list"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	StorageUnit string
	Value       int64
}

const (
	mebibyte = int64(1024 * 1024)
	gibibyte = 1024 * mebibyte
	tebibyte = 1024 * gibibyte

	// unlimitedStorage is the hard limit Harbor uses for an unlimited quota.
	unlimitedStorage = int64(-1)

	defaultStorageUnit = "GiB"
)

var storageUnits = []string{"MiB", "GiB", "TiB"}

func UpdateQuotaView(quta *models.Quota) string {
	var createView CreateView

	hard := quta.Hard["storage"]
	fmt.Printf("current storage: %v\n", currentStorage(hard))

	value, unit := storageDefaults(hard)
	createView.StorageUnit = unit

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

// currentStorage renders the hard limit the way the quota list command does,
// so an unlimited quota is not reported as a negative size.
func currentStorage(bytes int64) string {
	if bytes == unlimitedStorage {
		return "Unlimited"
	}

	return list.BytesToStorageString(bytes)
}

// storageDefaults converts a hard limit expressed in bytes into the value and
// unit used to prefill the form. The quota limit input only accepts whole
// numbers, so the largest unit that divides the limit evenly is picked. A
// limit that cannot be expressed that way, including an unlimited quota,
// yields an empty value so the user is asked for one instead of being handed a
// default the form itself rejects.
func storageDefaults(bytes int64) (value string, unit string) {
	if bytes == 0 {
		return "0", storageUnits[0]
	}

	if bytes > 0 {
		for _, u := range []struct {
			name string
			size int64
		}{
			{"TiB", tebibyte},
			{"GiB", gibibyte},
			{"MiB", mebibyte},
		} {
			if bytes >= u.size && bytes%u.size == 0 {
				return strconv.FormatInt(bytes/u.size, 10), u.name
			}
		}
	}

	return "", defaultStorageUnit
}
