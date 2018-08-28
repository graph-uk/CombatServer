package repositories

import (
	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/jinzhu/gorm"
)

// Cases repository
type Cases struct {
	context data.Context
}

//FindAll returns all cases from the database
func (t *Cases) FindAll() []models.Case {
	var cases []models.Case

	query := func(db *gorm.DB) {
		db.Find(&cases)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return cases
}

//FindBySessionID returns all cases for session from the database
func (t *Cases) FindBySessionID(sessionID string) []models.Case {
	var cases []models.Case

	query := func(db *gorm.DB) {
		db.Where(&models.Case{SessionID: sessionID}).Find(&cases)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return cases
}

// Find case by id
func (t *Cases) Find(id int) *models.Case {
	var result models.Case

	query := func(db *gorm.DB) {
		db.Find(&result, id)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return &result
}
