package models

// Try model
type Try struct {
	ID         int
	CaseID     int
	ExitStatus string `gorm:"size:50"`
	CaseOutput string `gorm:"Column:stdOut"`
}

// TableName ooverride
func (Try) TableName() string {
	return "profiles"
}
