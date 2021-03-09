package repositories

import (
	"math/rand"
	"time"

	"github.com/asdine/storm/q"

	"malibu-server/data"
	"malibu-server/data/models"
	"malibu-server/data/models/status"

	"github.com/asdine/storm"
)

// Cases repository
type Cases struct {
	context data.Context
}

// Create ...
func (t *Cases) Create(sessionCase *models.Case) error {
	query := func(db *storm.DB) {
		check(db.Save(sessionCase))
	}

	return t.context.Execute(query)
}

// Update ...
func (t *Cases) Update(sessionCase *models.Case) error {
	query := func(db *storm.DB) {
		check(db.Save(sessionCase))
	}

	return t.context.Execute(query)
}

//FindAll returns all cases from the database
func (t *Cases) FindAll() []models.Case {
	cases := &[]models.Case{}

	query := func(db *storm.DB) {
		check(db.All(cases))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return *cases
}

//FindBySessionID returns all cases for session from the database
func (t *Cases) FindBySessionID(sessionID string) []models.Case {
	cases := &[]models.Case{}

	query := func(db *storm.DB) {
		checkIgnore404(db.Find(`SessionID`, sessionID, cases))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return *cases
}

//FindProcessingCases returns cases with processing status
func (t *Cases) FindProcessingCases() []models.Case {
	cases := &[]models.Case{}

	query := func(db *storm.DB) {
		checkIgnore404(db.Find(`Status`, status.Processing, cases))
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}

	return *cases
}

// Find case by id
func (t *Cases) Find(id int) *models.Case {
	result := &models.Case{}

	query := func(db *storm.DB) {
		checkIgnore404(db.One(`ID`, id, result))
	}

	error := t.context.Execute(query)

	if error != nil || result.ID == 0 {
		return nil
	}

	return result
}

// StopCurrentCases marks all cases with Pending or Processing statuses to incomplete
func (t *Cases) StopCurrentCases() error {

	query := func(db *storm.DB) {
		cases := &[]models.Case{}
		sessionIDs := map[string]bool{}

		checkIgnore404(db.Select(q.Or(q.Eq(`Status`, status.Pending), q.Eq(`Status`, status.Processing))).Find(cases))

		for _, sessionCase := range *cases {
			sessionIDs[sessionCase.SessionID] = true
			sessionCase.Status = status.Incomplete
			check(db.Save(&sessionCase))
		}

		for sessionID := range sessionIDs {
			session := &models.Session{}
			check(db.One(`ID`, sessionID, session))
			session.Status = status.Incomplete
			check(db.Save(session))
		}
	}

	return t.context.Execute(query)
}

// AcquireFreeJob case by is not in progress and not finished
func (t *Cases) AcquireFreeJob() *models.Case {
	pendingJobs := []models.Case{}
	result := &models.Case{}
	session := &models.Session{}

	query := func(db *storm.DB) {
		checkIgnore404(db.Find(`Status`, status.Pending, &pendingJobs))
		if len(pendingJobs) > 0 {
			rand.Seed(time.Now().Unix())
			result = &pendingJobs[rand.Intn(len(pendingJobs))]

			result.Status = status.Processing
			result.DateStarted = time.Now()
			check(db.Save(result))

			check(db.One(`ID`, result.SessionID, session))
			session.Status = status.Processing
			check(db.Save(session))
		}
	}

	error := t.context.Execute(query)

	if error != nil || result.ID == 0 {
		return nil
	}

	return result
}

func (t *Cases) GetTotalCasesCountBySessionID(sessionID string) int {
	cases := &[]models.Case{}

	query := func(db *storm.DB) {
		check(db.Find(`SessionID`, sessionID, cases))
	}

	err := t.context.Execute(query)
	if err != nil {
		panic(err)
	}

	return len(*cases)
}

func (t *Cases) GetFailedCasesCountBySessionID(sessionID string) int {
	cases := &[]models.Case{}

	query := func(db *storm.DB) {
		checkIgnore404(db.Select(q.And(q.Eq(`SessionID`, sessionID), q.Eq(`Status`, status.Failed))).Find(cases))
	}

	err := t.context.Execute(query)
	if err != nil {
		panic(err)
	}

	return len(*cases)
}

func (t *Cases) DeleteByID(id int) {
	query := func(db *storm.DB) {
		check(db.Delete(`Case`, id))
	}

	t.context.Execute(query)
}
