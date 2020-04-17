package models

import (
	"time"

	"malibu-server/data/models/status"
)

// Session sdf
type Session struct {
	ID          string `storm:"id"`
	Arguments   string
	Status      status.Status
	Error       string
	DateCreated time.Time
}
