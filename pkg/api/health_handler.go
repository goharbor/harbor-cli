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
package api

import (
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/health"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func GetHealth() (*health.GetHealthOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize client context: ")
	}

	response, err := client.Health.GetHealth(ctx, &health.GetHealthParams{})
	if err != nil {
		switch err.(type) {
		case *health.GetHealthInternalServerError:
			return nil, fmt.Errorf("internal server error occurred while getting health status")
		default:
			return nil, fmt.Errorf("unknown error occurred while getting health status: %w", err)
		}
	}

	return response, nil
}
