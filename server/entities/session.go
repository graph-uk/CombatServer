package entities

import (
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//const (
//	session_progress    = 0
//	session_ok          = 1
//	session_falseNegErr = 2
//	session_appError    = 3
//)

type Session struct {
	ID                        uint   `gorm:"primary_key","AUTO_INCREMENT"`
	Params                    string `gorm:"type:string"` // params of tests name filtering, etc.
	HookFirstFail             bool   // sets to true, while first fail appeared, to prevent multiple slack alerting
	CasesExploringFailMessage string `gorm:"type:string"` // contains "combat.exe" error output, if combat cannot explore the session (for example - syntax error at any test).
}

func NewSession() *Session {
	var res Session
	//	res.ParentRun = t
	//	res.Run_ID = t.ID
	//	res.Worker_id = t.Worker_ID
	//	res.Status = session_progress
	//	res.StartTime = time.Now()
	//	db.Create(&res)

	//	res.loadLastOutFolder()
	//	exitCode := res.RunCombatTestBinary()
	//	res.saveTestOutFolder()

	//	if exitCode == 0 {
	//		res.Status = session_ok
	//		res.saveSuccessOut()
	//	} else {
	//		if res.isFalseNegFail() {
	//			res.Status = session_falseNegErr
	//		} else {
	//			res.Status = session_appError
	//			res.saveFailOut()
	//		}
	//	}

	//	res.SaveToDB(nil)
	return &res
}

func (t *Session) SaveToDB(tx *gorm.DB) {
	//	if tx == nil {
	//		db.Save(t)
	//	} else {
	//		tx.Save(t)
	//	}
}
