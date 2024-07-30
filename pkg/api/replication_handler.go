package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/replication"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func ListReplication(opts ...ListFlags) (*replication.ListReplicationPoliciesOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	var listFlags ListFlags

	if len(opts) > 0 {
		listFlags = opts[0]
	}

	response, err := client.Replication.ListReplicationPolicies(ctx, &replication.ListReplicationPoliciesParams{
		Page:     &listFlags.Page,
		PageSize: &listFlags.PageSize,
		Q:        &listFlags.Q,
		Name:     &listFlags.Name,
		Sort:     &listFlags.Sort,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
