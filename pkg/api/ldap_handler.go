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
