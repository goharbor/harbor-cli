package utils

import (
	"fmt"
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
