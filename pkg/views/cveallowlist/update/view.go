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

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type UpdateView struct {
	CveId      string
	IsExpire   bool
	ExpireDate string
}

func UpdateCveView(updateView *UpdateView) {
	theme := huh.ThemeCharm()
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("CVE ID").
				Value(&updateView.CveId).
				Description("CVE IDs are separator by commas").
				Validate(func(str string) error {
					if str == "" {
						return errors.New("cve id cannot be empty")
					}
					return nil
				}),
			huh.NewConfirm().
				Title("Expires").
				Value(&updateView.IsExpire).
				Affirmative("Date").
				Negative("never"),
		),
		huh.NewGroup(
			huh.NewInput().
				Validate(func(str string) error {
					if str == "" {
						return errors.New("ExpireDate cannot be empty")
					}
					return nil
				}).
				Description("Expire Date in the format YYYY/MM/DD").
				Title("Expire Date").
				Value(&updateView.ExpireDate),
		).WithHideFunc(func() bool {
			return !updateView.IsExpire
		}),
	).WithTheme(theme).Run()

	if err != nil {
		log.Fatal(err)
	}
}
