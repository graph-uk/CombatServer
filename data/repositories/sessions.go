package repositories

import (
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
		db.Find(&sessions)
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
