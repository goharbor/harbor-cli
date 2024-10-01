package api

import (
	"context"
	"fmt"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/permissions"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	
)

type PermissionsHandler struct {
	client *permissions.Client
}

func NewPermissionsHandler(client *permissions.Client) *PermissionsHandler {
	return &PermissionsHandler{
		client: client,
	}
}

func (h *PermissionsHandler) GetPermissions(ctx context.Context) (*models.Permissions, error) {
	params := &permissions.GetPermissionsParams{
		Context: ctx,
	}

	resp, err := h.client.GetPermissions(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	return resp.Payload, nil
}
