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

	"github.com/goharbor/harbor-cli/pkg/api"
)

// ValidateListFlags checks the pagination flags of the list commands.
func ValidateListFlags(opts api.ListFlags) error {
	if opts.Page < 1 {
		return fmt.Errorf("page number must be greater than or equal to 1")
	}
	if opts.PageSize < 0 {
		return fmt.Errorf("page size must be greater than or equal to 0")
	}
	if opts.PageSize > 100 {
		return fmt.Errorf("page size should be less than or equal to 100")
	}
	return nil
}
