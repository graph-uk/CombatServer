package models

import "time"

// Case model
type Case struct {
	ID          int    `gorm:"primary_key"`
	SessionID   string `gorm:"size:100"`
	Code        string `gorm:"size:100"`
	Title       string `gorm:"size:100"`
	CommandLine string `gorm:"size:500"`
	Status      int
	DateStarted time.Time

	Tries []Try `gorm:"foreignkey:CaseID"`
}
