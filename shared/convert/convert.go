package convert

import (
	"os"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// EnvOr returns the value of the environment variable key if set, otherwise fallback.
func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// NullableString converts an empty string to a nil pointer.
func NullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// FormatTimePtr formats a *time.Time into "2006-01-02 15:04:05" string.
// Returns empty string if t is nil or zero.
func FormatTimePtr(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// TimeToWrappers converts a *time.Time into a *wrapperspb.StringValue formatted as "2006-01-02 15:04:05".
// Returns nil if t is nil or zero.
func TimeToWrappers(t *time.Time) *wrapperspb.StringValue {
	if t == nil || t.IsZero() {
		return nil
	}
	return wrapperspb.String(t.Format("2006-01-02 15:04:05"))
}

// ParsePgTimestamp parses a Postgres timestamp string ("2006-01-02 15:04:05"
// or RFC3339) into a time.Time. Empty or unparsable input yields a zero time.
func ParsePgTimestamp(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return time.Time{}
		}
	}
	return t
}
