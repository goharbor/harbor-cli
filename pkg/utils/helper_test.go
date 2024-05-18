package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestFormatCreatedTime(t *testing.T) {
	// Test invalid timestamp
	invalidTimestamp := "invalid"
	_, err := FormatCreatedTime(invalidTimestamp)
	require.Error(t, err)

	// Test 2 minute ago
	minuteAgo := time.Now().Add(-2 * time.Minute).Format(time.RFC3339Nano)
	formatted, err := FormatCreatedTime(minuteAgo)
	require.NoError(t, err)
	require.Equal(t, "2 minute ago", formatted)

	// Test 3 hour ago
	hourAgo := time.Now().Add(-3 * time.Hour).Format(time.RFC3339Nano)
	formatted, err = FormatCreatedTime(hourAgo)
	require.NoError(t, err)
	require.Equal(t, "3 hour ago", formatted)

	// Test 5 day ago
	dayAgo := time.Now().Add(-24 * 5 * time.Hour).Format(time.RFC3339Nano)
	formatted, err = FormatCreatedTime(dayAgo)
	require.NoError(t, err)
	require.Equal(t, "5 day ago", formatted)
}
