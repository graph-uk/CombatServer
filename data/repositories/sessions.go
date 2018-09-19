package repositories

import (
	"strconv"
	"time"

	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/graph-uk/combat-server/data/models/status"
	"github.com/jinzhu/gorm"
)

// Sessions repository
type Sessions struct {
	context data.Context
}

//Create new session
func (t *Sessions) Create(arguments string, content []byte) *models.Session {
	session := &models.Session{
		ID:          strconv.FormatInt(time.Now().UnixNano(), 10),
		DateCreated: time.Now(),
		Status:      status.Awaiting,
		Arguments:   arguments}

	sessionFs := SessionsFS{}

	query := func(db *gorm.DB) {
		db.Create(session)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	go sessionFs.ProcessSession(session, content)

	return session
}

// Update ...
func (t *Sessions) Update(session *models.Session) error {
	query := func(db *gorm.DB) {
		db.Save(session)
	}

	return t.context.Execute(query)
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
