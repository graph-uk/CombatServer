package models

// Try model
type Try struct {
	ID         int `storm:"id,increment"`
	CaseID     int
	Output     string
	ExitStatus string

	//Items []string `gorm:"-"`
}

// TableName override
// func (Try) TableName() string {
// 	return "Tries"
// }
