package utils

import (
	"regexp"
	"strconv"
)

func ErrorStatusCode(err error) int64 {
	re := regexp.MustCompile(`\[(\d{3})\]`)
	matches := re.FindStringSubmatch(err.Error())
	var statusCode int64

	if len(matches) > 1 {
		stringCode := matches[1]
		statusCode, _ = strconv.ParseInt(stringCode, 10, 64)
	}
	return statusCode
}
