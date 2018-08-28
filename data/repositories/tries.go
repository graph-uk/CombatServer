package repositories

import (
	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/jinzhu/gorm"
)

// Tries repository
type Tries struct {
	context data.Context
}

//FindAll returns all tries from the database
func (t *Tries) FindAll() []models.Try {
	var tries []models.Try

	query := func(db *gorm.DB) {
		db.Find(&tries)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return tries
}

//FindByCaseID returns all tries for case from the database
func (t *Tries) FindByCaseID(caseID int) []models.Try {
	var tries []models.Try

	query := func(db *gorm.DB) {
		db.Where(&models.Try{CaseID: caseID}).Find(&tries)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return tries
}

// Find try by id
func (t *Tries) Find(id int) *models.Try {
	var try models.Try

	query := func(db *gorm.DB) {
		db.Find(&try, id)
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return &try
}
