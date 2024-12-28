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
package views

import (
	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

func ConfirmElevation() (bool, error) {
	var confirm bool

	err := huh.NewConfirm().
		Title("Are you sure to elevate the user to admin role?").
		Affirmative("Yes").
		Negative("No").
		Value(&confirm).Run()
	if err != nil {
		log.Fatal(err)
	}

	return confirm, nil
}
