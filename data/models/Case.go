package models

import "time"

// Case model
type Case struct {
	ID           int
	CommandLine  string    `gorm:"Column:cmdLine"`
	SessionID    string    `gorm:"Column:sessionID"`
	IsInProgress bool      `gorm:"Column:inProgress"`
	IsFinished   bool      `gorm:"Column:finished"`
	IsPassed     bool      `gorm:"Column:passed"`
	DateStarted  time.Time `gorm:"Column:startedAt"`
	Tries        []Try     `gorm:"foreignkey:CaseID"`
}
