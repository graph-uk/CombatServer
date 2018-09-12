package entities

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Case struct {
	ID         uint   `gorm:"primary_key","AUTO_INCREMENT"`
	CmdLine    string `gorm:"type:string"`
	SessionID  string `gorm:"type:string"`
	InProgress bool   `gorm:"default:false"`
	Finished   bool   `gorm:"default:false"`
	Passed     bool   `gorm:"default:false"`
	StartedAt  time.Time
}

func NewCase() *Case {
	var res Case

	res.SaveToDB(nil)
	return &res
}

func (t *Case) SaveToDB(tx *gorm.DB) {
	if tx == nil {
		tx.Save(t)
	} else {
		tx.Save(t)
	}
}
