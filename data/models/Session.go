package models

import "time"

// Session sdf
type Session struct {
	ID          string `gorm:"primary_key,type:varchar(20)"`
	Arguments   string `gorm:"size:500"`
	Status      int
	DateCreated time.Time
}
