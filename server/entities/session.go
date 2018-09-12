package entities

import (
	"github.com/jinzhu/gorm"
)

type Session struct {
	ID                        uint   `gorm:"primary_key","AUTO_INCREMENT"`
	Params                    string `gorm:"type:string"`   // params of tests name filtering, etc.
	HookFirstFail             bool   `gorm:"default:false"` // sets to true, while first fail appeared, to prevent multiple slack alerting
	CasesExploringFailMessage string `gorm:"type:string"`   // contains "combat.exe" error output, if combat cannot explore the session (for example - syntax error at any test).
}

func NewSession() *Session {
	var res Session

	res.SaveToDB(nil)
	return &res
}

func (t *Session) SaveToDB(tx *gorm.DB) {
	if tx == nil {
		tx.Save(t)
	} else {
		tx.Save(t)
	}
}
