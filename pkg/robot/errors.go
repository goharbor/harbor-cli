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
package robot

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/utils"
)

// WrapError turns a Harbor API error into the message the robot commands
// show, with a uniform hint for missing admin privileges.
func WrapError(action string, err error) error {
	if utils.ParseHarborErrorCode(err) == "403" {
		return fmt.Errorf("permission denied: admin privileges are required to %s", action)
	}
	return fmt.Errorf("failed to %s: %v", action, utils.ParseHarborErrorMsg(err))
}
