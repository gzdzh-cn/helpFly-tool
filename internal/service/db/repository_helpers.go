package db

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

func isNoRows(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func parseDBTime(value string) time.Time {
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
	} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

func formatDBTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
