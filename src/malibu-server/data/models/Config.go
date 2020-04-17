package models

import (
	"time"
)

// Config model
type Config struct {
	ID                  int `storm:"id,increment"`
	MuteTimestamp       time.Time
	NotificationEnabled bool
}
