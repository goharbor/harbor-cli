package api

import (
	"github.com/goharbor/go-client/pkg/sdk/v2.0/client/ldap"
	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func LdapSearchUser(username string) (*ldap.SearchLdapUserOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	res, err := client.Ldap.SearchLdapUser(ctx, &ldap.SearchLdapUserParams{
		Username: &username,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func LdapPingServer(ldapConf *models.LdapConf) (*ldap.PingLdapOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	res, err := client.Ldap.PingLdap(ctx, &ldap.PingLdapParams{
		Ldapconf: ldapConf,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
