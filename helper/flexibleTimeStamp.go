package helper

import (
	"fmt"
	"time"
)

func ParseFlexibleTimestamp(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid timestamp: %s", value)
}
