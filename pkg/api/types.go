package api

import "github.com/goharbor/go-client/pkg/sdk/v2.0/models"

type ListFlags struct {
	Name     string
	Page     int64
	PageSize int64
	Q        string
	Sort     string
	Public   bool
}

// Provides type for List member
type ListMemberOptions struct {
	ProjectNameOrID string
	Page            int64
	PageSize        int64
	EntityName      string
	WithDetail      bool
}

// Provides type for Update member
type UpdateMemberOptions struct {
	ID              int64
	ProjectNameOrID string
	RoleID          *models.RoleRequest
}

// Provides Params for getting the Member
type GetMemberOptions struct {
	ID              int64
	ProjectNameOrID string
}
