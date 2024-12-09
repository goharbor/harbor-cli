package errors

import (
	"github.com/goharbor/harbor-cli/pkg/utils"
)

func getErrorMessage(statusCode int64, action string) string {
	switch statusCode {
	case 400:
		if action == "create" {
			return "Give adequate information to create a registry."
		}
	case 401:
		return "Registry Unauthorized"
	case 403:
		return "Registry Forbidden"
	case 404:
		return "Registry is not found"
	case 409:
		if action == "create" {
			return "Registries with the same name exist"
		}
		return "Conflict"
	case 412:
		return "Registry is not found"
	case 500:
		return "Internal Server Error"
	}
	return "cannot identify the status code"
}

func ErrorCreateRegistry(err error) string {
	return getErrorMessage(utils.ErrorStatusCode(err), "create")
}

func ErrorDeleteRegistry(err error) string {
	return getErrorMessage(utils.ErrorStatusCode(err), "delete")
}

func ErrorViewRegistry(err error) string {
	return getErrorMessage(utils.ErrorStatusCode(err), "view")
}

func ErrorListRegistry(err error) string {
	return getErrorMessage(utils.ErrorStatusCode(err), "list")
}

func ErrorUpdateRegistry(err error) string {
	return getErrorMessage(utils.ErrorStatusCode(err), "update")
}
