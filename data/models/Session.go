package models

import "time"

// Session sdf
type Session struct {
	ID          string `gorm:"type:varchar(20)"`
	Arguments   string `gorm:"size:500"`
	Status      int
	DateCreated time.Time
}
