package models

import (
	"time"
)

// Config model
type Config struct {
	ID                  int `gorm:"primary_key"`
	MuteTimestamp       time.Time
	NotificationEnabled bool
}
