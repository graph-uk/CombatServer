package models

// Try model
type Try struct {
	ID         int `storm:"id,increment"`
	CaseID     int
	Output     string
	ExitStatus string
}
