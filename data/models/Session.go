package models

// Session sdf
type Session struct {
	ID                        string `gorm:"type:varchar(20)"`
	Params                    string `gorm:"size:50"`
	CasesExploringFailMessage string
	IsFirstFail               bool `gorm:"Column:hook_FirstFail"`
}
