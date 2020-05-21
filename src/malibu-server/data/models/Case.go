package models

import (
	"crypto/md5"
	"encoding/hex"
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
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//identity used for storing last successfull tries of the case, across all sessions
func (t *Case) GetCmdHash() string {
	return GetMD5Hash(t.CommandLine)
}
