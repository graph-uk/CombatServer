package repositories

import (
	"fmt"
	"time"

	"malibu-server/data"
	"malibu-server/data/models"
	"malibu-server/data/models/status"

	"github.com/jinzhu/gorm"
)

// Cases repository
type Cases struct {
	context data.Context
}

// Create ...
func (t *Cases) Create(sessionCase *models.Case) error {
	query := func(db *gorm.DB) {
		db.Create(sessionCase)
	}

	return t.context.Execute(query)
}

// Update ...
func (t *Cases) Update(sessionCase *models.Case) error {
	query := func(db *gorm.DB) {
		db.Save(sessionCase)
	}

	return t.context.Execute(query)
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

//FindProcessingCases returns cases with processing status
func (t *Cases) FindProcessingCases() []models.Case {
	var cases []models.Case

	query := func(db *gorm.DB) {
		db.Where(&models.Case{Status: status.Processing}).Find(&cases)
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

// StopCurrentCases marks all cases with Pending or Processing statuses to incomplete
func (t *Cases) StopCurrentCases() error {

	query := func(db *gorm.DB) {
		var cases []models.Case
		sessionIDs := map[string]bool{}

		db.Where(&models.Case{Status: status.Pending}).Or(&models.Case{Status: status.Processing}).Find(&cases)

		for _, sessionCase := range cases {
			fmt.Println(sessionCase.ID)
			sessionIDs[sessionCase.SessionID] = true
			sessionCase.Status = status.Incomplete
			db.Save(&sessionCase)
		}

		for sessionID := range sessionIDs {
			var session models.Session
			db.Find(&session, sessionID)
			session.Status = status.Incomplete
			db.Save(session)
		}
	}

	return t.context.Execute(query)
}

// AcquireFreeJob case by is not in progress and not finished
func (t *Cases) AcquireFreeJob() *models.Case {
	var result models.Case
	var session models.Session

	query := func(db *gorm.DB) {
		// Where is string because of shitty gorm which can't filter by false :-(
		db.Order("random()").Where(&models.Case{Status: status.Pending}).First(&result)
		if result.ID > 0 {
			result.Status = status.Processing
			result.DateStarted = time.Now()
			db.Save(&result)

			db.Find(&session, result.SessionID)
			session.Status = status.Processing
			db.Save(&session)
		}
	}

	error := t.context.Execute(query)

	if error != nil || result.ID == 0 {
		return nil
	}

	return &result
}

func (t *Cases) GetTotalCasesCountBySessionID(sessionID string) int {
	var cases []models.Case

	query := func(db *gorm.DB) {
		db.Where(&models.Case{SessionID: sessionID}).Find(&cases)
	}

	err := t.context.Execute(query)
	if err != nil {
		panic(err)
	}

	return len(cases)
}

func (t *Cases) GetFailedCasesCountBySessionID(sessionID string) int {
	var cases []models.Case

	query := func(db *gorm.DB) {
		db.Where(&models.Case{SessionID: sessionID, Status: status.Failed}).Find(&cases)
	}

	err := t.context.Execute(query)
	if err != nil {
		panic(err)
	}

	return len(cases)
}
