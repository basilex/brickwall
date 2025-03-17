package common

import (
	"time"
)

func TimestampISO3339NS() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999999999+00:00")
}
