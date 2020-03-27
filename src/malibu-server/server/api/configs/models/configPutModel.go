package configs

import (
	"time"
)

type ConfigPutModel struct {
	MuteTimestamp       time.Time
	NotificationEnabled bool
}
