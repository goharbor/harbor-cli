package errors

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type harborErrorPayload struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func isHarborError(err error) *Error {
	var e *Error
	if as(err, &e) {
		return e
	}
	return nil
}

func parseHarborErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	val := reflect.ValueOf(err)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	field := val.FieldByName("Payload")
	if field.IsValid() {
		payload := field.Interface()
		jsonBytes, jsonErr := json.Marshal(payload)
		if jsonErr == nil {
			var harborErr harborErrorPayload
			if unmarshalErr := json.Unmarshal(jsonBytes, &harborErr); unmarshalErr == nil {
				if len(harborErr.Errors) > 0 {
					return harborErr.Errors[0].Message
				}
			}
		}
	}
	return fmt.Sprintf("%v", err.Error())
}

func parseHarborErrorCode(err error) string {
	parts := strings.Split(err.Error(), "]")
	if len(parts) >= 2 {
		codePart := strings.TrimSpace(parts[1])
		if strings.HasPrefix(codePart, "[") && len(codePart) == 4 {
			code := codePart[1:4]
			return code
		}
	}
	return ""
}
