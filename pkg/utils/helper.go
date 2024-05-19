package utils

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

func FormatCreatedTime(timestamp string) (string, error) {
	t, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return "", err
	}

	duration := time.Since(t)

	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := int(duration.Hours() / 24)

	if minutes < 60 {
		return fmt.Sprintf("%d minute ago", minutes), nil
	} else if hours < 24 {
		return fmt.Sprintf("%d hour ago", hours), nil
	} else {
		return fmt.Sprintf("%d day ago", days), nil
	}
}

func HandleURLFormat(server string) (string, error) {
	server = strings.TrimSpace(server)
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "https://" + server
	}
	parsedURL, err := url.Parse(server)
	if err != nil {
		return "", fmt.Errorf("invalid server address: %s", err)
	}
	if parsedURL.Port() == "" {
		if parsedURL.Scheme == "https" {
			parsedURL.Host = parsedURL.Host + ":443"
		} else {
			parsedURL.Host = parsedURL.Host + ":80"
		}
	}
	return parsedURL.String(), nil
}
