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
package create

import (
	"errors"
	"strconv"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

type CreateView struct {
	ScopeSelectors   RetentionSelector
	TagSelectors     RetentionSelector
	KeepLatestPushed int64
	Cron             string
}

type RetentionSelector struct {
	Decoration string
	Pattern    string
}

func CreateRetentionView(createView *CreateView) {
	keepLatestRaw := "10"
	if createView.KeepLatestPushed == 0 {
		createView.KeepLatestPushed = 10
	} else {
		keepLatestRaw = strconv.FormatInt(createView.KeepLatestPushed, 10)
	}
	if createView.Cron == "" {
		createView.Cron = "0 0 0 * * *"
	}

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("For the repositories").
				Options(
					huh.NewOption("matching", "repoMatches"),
					huh.NewOption("excluding", "repoExcludes"),
				).
				Value(&createView.ScopeSelectors.Decoration).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("repository decoration cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Repository selector pattern").
				Description("Enter **, repo*, or comma-separated patterns").
				Value(&createView.ScopeSelectors.Pattern).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("repository pattern cannot be empty")
					}
					return nil
				}),
			huh.NewSelect[string]().
				Title("For the tags").
				Options(
					huh.NewOption("matching", "matches"),
					huh.NewOption("excluding", "excludes"),
				).
				Value(&createView.TagSelectors.Decoration).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("tag decoration cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Tag selector pattern").
				Description("Enter **, v*, or comma-separated patterns").
				Value(&createView.TagSelectors.Pattern).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("tag pattern cannot be empty")
					}
					return nil
				}),
			huh.NewInput().
				Title("Keep latest pushed artifacts").
				Description("Default 10").
				Value(&keepLatestRaw).
				Validate(func(v string) error {
					n, err := strconv.ParseInt(v, 10, 64)
					if err != nil || n <= 0 {
						return errors.New("keep latest value must be greater than 0")
					}
					return nil
				}),
			huh.NewInput().
				Title("Cron schedule").
				Description("Default 0 0 0 * * *").
				Value(&createView.Cron).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("cron cannot be empty")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).Run()

	if err != nil {
		log.Fatal(err)
	}

	parsed, err := strconv.ParseInt(keepLatestRaw, 10, 64)
	if err != nil || parsed <= 0 {
		createView.KeepLatestPushed = 10
		return
	}

	createView.KeepLatestPushed = parsed
}
