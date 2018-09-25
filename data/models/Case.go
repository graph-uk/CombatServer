package models

import (
	"time"

	"github.com/graph-uk/combat-server/data/models/status"
)

// Case model
type Case struct {
	ID          int    `gorm:"primary_key"`
	SessionID   string `gorm:"size:100"`
	Code        string `gorm:"size:100"`
	Title       string `gorm:"size:100"`
	CommandLine string `gorm:"size:500"`
	Status      status.Status
	DateStarted time.Time

	Tries []Try `gorm:"foreignkey:CaseID"`
}
