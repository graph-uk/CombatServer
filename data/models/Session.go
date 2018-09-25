package models

import (
	"time"

	"github.com/graph-uk/combat-server/data/models/status"
)

// Session sdf
type Session struct {
	ID          string `gorm:"primary_key,type:varchar(20)"`
	Arguments   string `gorm:"size:500"`
	Status      status.Status
	Error       string `gorm:"size:1000"`
	DateCreated time.Time
}
