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

func LdapImportUser(uids []string) error {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return err
	}

	_, err = client.Ldap.ImportLdapUser(ctx, &ldap.ImportLdapUserParams{
		UIDList: &models.LdapImportUsers{
			LdapUIDList: uids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func LdapSearchGroup(groupName, groupDN string) (*ldap.SearchLdapGroupOK, error) {
	ctx, client, err := utils.ContextWithClient()
	if err != nil {
		return nil, err
	}

	res, err := client.Ldap.SearchLdapGroup(ctx, &ldap.SearchLdapGroupParams{
		Groupdn:   &groupDN,
		Groupname: &groupName,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
