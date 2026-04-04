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
package schedules

import (
	"fmt"

	"github.com/goharbor/harbor-cli/pkg/utils"
)

func formatScheduleError(operation string, err error, requiredPermission string) error {
	errorCode := utils.ParseHarborErrorCode(err)

	switch errorCode {
	case "400":
		return fmt.Errorf("%s: invalid request. For schedule status use job_type=all; for queue action use stop|pause|resume", operation)
	case "401":
		return fmt.Errorf("%s: authentication required. Please run 'harbor login' and try again", operation)
	case "403":
		if requiredPermission == "authenticated" {
			return fmt.Errorf("%s: permission denied. Your account is authenticated but lacks access", operation)
		}
		return fmt.Errorf("%s: permission denied. This operation requires %s on jobservice-monitor", operation, requiredPermission)
	case "404":
		return fmt.Errorf("%s: resource not found or not accessible in current context", operation)
	case "422":
		return fmt.Errorf("%s: request validation failed. Please check request body and action values", operation)
	case "500":
		return fmt.Errorf("%s: Harbor internal error. Retry and check Harbor server logs", operation)
	default:
		msg := utils.ParseHarborErrorMsg(err)
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("%s: %s", operation, msg)
	}
}
