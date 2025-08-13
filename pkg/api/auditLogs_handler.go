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
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/auditlog"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func AuditLogs(opts ListFlags) (*auditlog.ListAuditLogExtsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Auditlog.ListAuditLogExts(ctx,
		&auditlog.ListAuditLogExtsParams{
			Q:        &opts.Q,
			Sort:     &opts.Sort,
			Page:     &opts.Page,
			PageSize: &opts.PageSize,
		})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func AuditLogsLegacy(opts ListFlags) (*auditlog.ListAuditLogsOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Auditlog.ListAuditLogs(ctx,
		&auditlog.ListAuditLogsParams{
			Q:        &opts.Q,
			Sort:     &opts.Sort,
			Page:     &opts.Page,
			PageSize: &opts.PageSize,
		})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func AuditLogEventTypes() (*auditlog.ListAuditLogEventTypesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Auditlog.ListAuditLogEventTypes(ctx,
		&auditlog.ListAuditLogEventTypesParams{})
	if err != nil {
		return nil, err
	}

	return response, nil
}
