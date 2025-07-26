package utils

import (
	"reflect"

	"github.com/goharbor/go-client/pkg/sdk/v2.0/models"
)

// ConfigurationCategories maps each configuration field to its category
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

	// Security - Security policies, certificates, permissions
	"LdapVerifyCert":             "security",
	"OIDCVerifyCert":             "security",
	"UaaVerifyCert":              "security",
	"HTTPAuthproxyVerifyCert":    "security",
	"SelfRegistration":           "security",
	"ProjectCreationRestriction": "security",
	"RobotTokenDuration":         "security",
	"TokenExpiration":            "security",
	"SessionTimeout":             "security",
	"RobotNamePrefix":            "security",

	// System - General system behavior, storage, auditing
	"ReadOnly":                   "system",
	"QuotaPerProjectEnable":      "system",
	"StoragePerProject":          "system",
	"NotificationEnable":         "system",
	"ScannerSkipUpdatePulltime":  "system",
	"SkipAuditLogDatabase":       "system",
	"AuditLogForwardEndpoint":    "system",
	"BannerMessage":              "system",
	"DisabledAuditLogEventTypes": "system",
	"OIDCLogout":                 "system",
}

// GetValidCategories returns all available categories
func GetValidCategories() []string {
	categories := make(map[string]bool)
	for _, category := range ConfigurationCategories {
		categories[category] = true
	}

	result := make([]string, 0, len(categories))
	for category := range categories {
		result = append(result, category)
	}
	return result
}

// GetConfigurationsByCategory filters configurations by category
func GetConfigurationsByCategory(configs *models.Configurations, category string) map[string]interface{} {
	result := make(map[string]interface{})

	configValue := reflect.ValueOf(configs).Elem()
	configType := configValue.Type()

	for i := 0; i < configValue.NumField(); i++ {
		fieldName := configType.Field(i).Name
		fieldValue := configValue.Field(i)

		// Check if this field belongs to the requested category
		if ConfigurationCategories[fieldName] == category {
			if fieldValue.IsValid() && !fieldValue.IsNil() {
				result[fieldName] = fieldValue.Interface()
			}
		}
	}

	return result
}

// GetAllCategorizedConfigurations returns configurations grouped by category
func GetAllCategorizedConfigurations(configs *models.Configurations) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})

	// Initialize categories
	for _, category := range GetValidCategories() {
		result[category] = make(map[string]interface{})
	}

	configValue := reflect.ValueOf(configs).Elem()
	configType := configValue.Type()

	for i := 0; i < configValue.NumField(); i++ {
		fieldName := configType.Field(i).Name
		fieldValue := configValue.Field(i)

		if category, exists := ConfigurationCategories[fieldName]; exists {
			if fieldValue.IsValid() && !fieldValue.IsNil() {
				result[category][fieldName] = fieldValue.Interface()
			}
		}
	}

	return result
}
