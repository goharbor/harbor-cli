package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type HarborErrorPayload struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func ParseHarborError(err error) string {
	if err == nil {
		return ""
	}

	val := reflect.ValueOf(err)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field := val.FieldByName("Payload")
	if field.IsValid() {
		payload := field.Interface()
		jsonBytes, jsonErr := json.Marshal(payload)
		if jsonErr == nil {
			var harborErr HarborErrorPayload
			if unmarshalErr := json.Unmarshal(jsonBytes, &harborErr); unmarshalErr == nil {
				if len(harborErr.Errors) > 0 {
					return fmt.Sprintf("%s field", harborErr.Errors[0].Message)
				}
			}
		}
	}
	return fmt.Sprintf("Error: %s", err.Error())
}
