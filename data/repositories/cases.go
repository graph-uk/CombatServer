package repositories

import (
	"time"

	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/graph-uk/combat-server/data/models/status"
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

// AcquireFreeJob case by is not in progress and not finished
func (t *Cases) AcquireFreeJob() *models.Case {
	var result models.Case

	query := func(db *gorm.DB) {
		// Where is string because of shitty gorm which can't filter by false :-(
		db.Where("finished = 0 AND inProgress = 0").First(&result)
		if result.ID > 0 {
			result.Status = status.Awaiting
			result.DateStarted = time.Now()
			db.Save(&result)
		}
	}

	error := t.context.Execute(query)

	if error != nil || result.ID == 0 {
		return nil
	}

	return &result
}
