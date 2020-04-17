package models

import (
	"time"

	"malibu-server/data/models/status"
)

// Case model
type Case struct {
	ID          int `storm:"id,increment"`
	SessionID   string
	Code        string 
	Title       string 
	CommandLine string 
	Status      status.Status
	DateStarted time.Time

	//Tries []Try 
}
