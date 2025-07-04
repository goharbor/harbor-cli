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
