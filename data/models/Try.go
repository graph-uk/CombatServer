package models

// Try model
type Try struct {
	ID         int `gorm:"primary_key"`
	CaseID     int
	Output     string `gorm:"size:1000"`
	ExitStatus string `gorm:"size:50"`

	Items []string `gorm:"-"`
}

// TableName override
func (Try) TableName() string {
	return "Tries"
}
