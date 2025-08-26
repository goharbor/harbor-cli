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
package utils

var ConfigurationCategories = map[string]string{
	// Authentication - User/service authentication methods
	"AuthMode":                         "authentication",
	"LdapBaseDn":                       "authentication",
	"LdapFilter":                       "authentication",
	"LdapGroupAdminDn":                 "authentication",
	"LdapGroupAttributeName":           "authentication",
	"LdapGroupBaseDn":                  "authentication",
	"LdapGroupMembershipAttribute":     "authentication",
	"LdapGroupSearchFilter":            "authentication",
	"LdapGroupSearchScope":             "authentication",
	"LdapScope":                        "authentication",
	"LdapSearchDn":                     "authentication",
	"LdapSearchPassword":               "authentication",
	"LdapTimeout":                      "authentication",
	"LdapUID":                          "authentication",
	"LdapURL":                          "authentication",
	"OIDCAdminGroup":                   "authentication",
	"OIDCAutoOnboard":                  "authentication",
	"OIDCClientID":                     "authentication",
	"OIDCClientSecret":                 "authentication",
	"OIDCEndpoint":                     "authentication",
	"OIDCExtraRedirectParms":           "authentication",
	"OIDCGroupFilter":                  "authentication",
	"OIDCGroupsClaim":                  "authentication",
	"OIDCName":                         "authentication",
	"OIDCScope":                        "authentication",
	"OIDCUserClaim":                    "authentication",
	"HTTPAuthproxyAdminGroups":         "authentication",
	"HTTPAuthproxyAdminUsernames":      "authentication",
	"HTTPAuthproxyEndpoint":            "authentication",
	"HTTPAuthproxyServerCertificate":   "authentication",
	"HTTPAuthproxySkipSearch":          "authentication",
	"HTTPAuthproxyTokenreviewEndpoint": "authentication",
	"UaaClientID":                      "authentication",
	"UaaClientSecret":                  "authentication",
	"UaaEndpoint":                      "authentication",
	"PrimaryAuthMode":                  "authentication",
	"LdapVerifyCert":                   "authentication",
	"OIDCVerifyCert":                   "authentication",
	"OIDCLogout":                       "authentication",

	// Security - Security policies, certificates, permissions
	"UaaVerifyCert":              "security",
	"HTTPAuthproxyVerifyCert":    "security",
	"SelfRegistration":           "security",
	"ProjectCreationRestriction": "security",
	"RobotTokenDuration":         "security",
	"TokenExpiration":            "security",
	"SessionTimeout":             "security",

	// System - General system behavior, storage, auditing
	"RobotNamePrefix":            "system",
	"ReadOnly":                   "system",
	"QuotaPerProjectEnable":      "system",
	"StoragePerProject":          "system",
	"NotificationEnable":         "system",
	"ScannerSkipUpdatePulltime":  "system",
	"SkipAuditLogDatabase":       "system",
	"AuditLogForwardEndpoint":    "system",
	"BannerMessage":              "system",
	"DisabledAuditLogEventTypes": "system",
}

func IsCategory(fieldName string, category string) bool {
	if category == "" {
		return true // If no category is specified, return true for all fields
	}
	if cat, exists := ConfigurationCategories[fieldName]; exists {
		return cat == category
	}
	return false
}
