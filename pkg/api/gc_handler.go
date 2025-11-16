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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/gc"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListGCs(opts *ListFlags) (gc.GetGCHistoryOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return gc.GetGCHistoryOK{}, err
	}

	response, err := client.GC.GetGCHistory(ctx, &gc.GetGCHistoryParams{
		Q:        &opts.Q,
		Sort:     &opts.Sort,
		Page:     &opts.Page,
		PageSize: &opts.PageSize,
	})
	if err != nil {
		return gc.GetGCHistoryOK{}, err
	}

	return *response, nil
}

func StopGC(id int64) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.GC.StopGC(ctx, &gc.StopGCParams{
		GCID: id,
	})
	if err != nil {
		return err
	}

	return nil
}
