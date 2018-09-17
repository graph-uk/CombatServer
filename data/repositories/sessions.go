package repositories

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/jinzhu/gorm"
)

// Sessions repository
type Sessions struct {
	context data.Context
}

//FindAll returns all sessions from the database
func (t *Sessions) FindAll() []models.Session {
	var sessions []models.Session

	query := func(db *gorm.DB) {
		db.Order("id desc").Find(&sessions)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return sessions
}

// Find session by id
func (t *Sessions) Find(id string) *models.Session {
	var session models.Session

	query := func(db *gorm.DB) {
		db.Find(&session, id)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return &session
}

// FindSessionContent returns session archive in BASE64 format from local disk
func (t *Sessions) FindSessionContent(sessionID string) string {
	zipFile, err := ioutil.ReadFile("./sessions/" + sessionID + "/archived.zip")

	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(zipFile)
}
