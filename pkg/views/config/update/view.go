package update

import "github.com/goharbor/go-client/pkg/sdk/v2.0/models"

// Function to convert ConfigurationsResponse to Configurations
func ConvertToConfigurations(
	resp *models.ConfigurationsResponse,
	ldapSearchPassword string,
	oidcClientSecret string,
) *models.Configurations {
	return &models.Configurations{
		AuditLogForwardEndpoint:          getStringValue(resp.AuditLogForwardEndpoint),
		AuthMode:                         getStringValue(resp.AuthMode),
		BannerMessage:                    getStringValue(resp.BannerMessage),
		HTTPAuthproxyAdminGroups:         getStringValue(resp.HTTPAuthproxyAdminGroups),
		HTTPAuthproxyAdminUsernames:      getStringValue(resp.HTTPAuthproxyAdminUsernames),
		HTTPAuthproxyEndpoint:            getStringValue(resp.HTTPAuthproxyEndpoint),
		HTTPAuthproxyServerCertificate:   getStringValue(resp.HTTPAuthproxyServerCertificate),
		HTTPAuthproxySkipSearch:          getBoolValue(resp.HTTPAuthproxySkipSearch),
		HTTPAuthproxyTokenreviewEndpoint: getStringValue(resp.HTTPAuthproxyTokenreviewEndpoint),
		HTTPAuthproxyVerifyCert:          getBoolValue(resp.HTTPAuthproxyVerifyCert),
		LdapBaseDn:                       getStringValue(resp.LdapBaseDn),
		LdapFilter:                       getStringValue(resp.LdapFilter),
		LdapGroupAdminDn:                 getStringValue(resp.LdapGroupAdminDn),
		LdapGroupAttributeName:           getStringValue(resp.LdapGroupAttributeName),
		LdapGroupBaseDn:                  getStringValue(resp.LdapGroupBaseDn),
		LdapGroupMembershipAttribute:     getStringValue(resp.LdapGroupMembershipAttribute),
		LdapGroupSearchFilter:            getStringValue(resp.LdapGroupSearchFilter),
		LdapGroupSearchScope:             getInt64Value(resp.LdapGroupSearchScope),
		LdapScope:                        getInt64Value(resp.LdapScope),
		LdapSearchDn:                     getStringValue(resp.LdapSearchDn),
		LdapSearchPassword:               &ldapSearchPassword,
		LdapTimeout:                      getInt64Value(resp.LdapTimeout),
		LdapUID:                          getStringValue(resp.LdapUID),
		LdapURL:                          getStringValue(resp.LdapURL),
		LdapVerifyCert:                   getBoolValue(resp.LdapVerifyCert),
		NotificationEnable:               getBoolValue(resp.NotificationEnable),
		OIDCAdminGroup:                   getStringValue(resp.OIDCAdminGroup),
		OIDCAutoOnboard:                  getBoolValue(resp.OIDCAutoOnboard),
		OIDCClientID:                     getStringValue(resp.OIDCClientID),
		OIDCClientSecret:                 &oidcClientSecret,
		OIDCEndpoint:                     getStringValue(resp.OIDCEndpoint),
		OIDCExtraRedirectParms:           getStringValue(resp.OIDCExtraRedirectParms),
		OIDCGroupFilter:                  getStringValue(resp.OIDCGroupFilter),
		OIDCGroupsClaim:                  getStringValue(resp.OIDCGroupsClaim),
		OIDCName:                         getStringValue(resp.OIDCName),
		OIDCScope:                        getStringValue(resp.OIDCScope),
		OIDCUserClaim:                    getStringValue(resp.OIDCUserClaim),
		OIDCVerifyCert:                   getBoolValue(resp.OIDCVerifyCert),
		PrimaryAuthMode:                  getBoolValue(resp.PrimaryAuthMode),
		ProjectCreationRestriction:       getStringValue(resp.ProjectCreationRestriction),
		QuotaPerProjectEnable:            getBoolValue(resp.QuotaPerProjectEnable),
		ReadOnly:                         getBoolValue(resp.ReadOnly),
		RobotNamePrefix:                  getStringValue(resp.RobotNamePrefix),
		RobotTokenDuration:               getInt64Value(resp.RobotTokenDuration),
		ScannerSkipUpdatePulltime:        getBoolValue(resp.ScannerSkipUpdatePulltime),
		SelfRegistration:                 getBoolValue(resp.SelfRegistration),
		SessionTimeout:                   getInt64Value(resp.SessionTimeout),
		SkipAuditLogDatabase:             getBoolValue(resp.SkipAuditLogDatabase),
		StoragePerProject:                getInt64Value(resp.StoragePerProject),
		TokenExpiration:                  getInt64Value(resp.TokenExpiration),
		UaaClientID:                      getStringValue(resp.UaaClientID),
		UaaClientSecret:                  getStringValue(resp.UaaClientSecret),
		UaaEndpoint:                      getStringValue(resp.UaaEndpoint),
		UaaVerifyCert:                    getBoolValue(resp.UaaVerifyCert),
	}
}

// Helper functions to extract values from configuration items
func getStringValue(item *models.StringConfigItem) *string {
	if item != nil {
		return &item.Value
	}
	return nil
}

func getBoolValue(item *models.BoolConfigItem) *bool {
	if item != nil {
		return &item.Value
	}
	return nil
}

func getInt64Value(item *models.IntegerConfigItem) *int64 {
	if item != nil {
		return &item.Value
	}
	return nil
}
