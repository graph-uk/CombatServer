package configs

import (
	"time"
)

type ConfigModel struct {
	MuteTimestamp       time.Time
	NotificationEnabled bool
	MuteDurationMinutes int
}
