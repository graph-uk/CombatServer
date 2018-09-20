package repositories

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/graph-uk/combat-server/data"
	"github.com/graph-uk/combat-server/data/models"
	"github.com/graph-uk/combat-server/data/models/status"
	"github.com/graph-uk/combat-server/server/config"
	"github.com/jinzhu/gorm"
	"github.com/mholt/archiver"
)

const triesPathTemplate = "_data/tries/%d"
const triesArchivePathTemplate = "_data/tries/%d/archived.zip"
const triesUnarchivedPathTemplate = "_data/tries/%d/_"

// Tries repository
type Tries struct {
	context data.Context
}

// Create Try
func (t *Tries) Create(try *models.Try, content []byte) error {
	var session models.Session
	var sessionCase models.Case

	query := func(db *gorm.DB) {
		db.Find(&sessionCase, try.CaseID)
		db.Find(&session, sessionCase.SessionID)

		db.Create(try)
	}

	err := t.context.Execute(query)

	if err != nil {
		return err
	}

	path := fmt.Sprintf(triesPathTemplate, try.ID)
	archivedPath := fmt.Sprintf(triesArchivePathTemplate, try.ID)
	unarchivedPath := fmt.Sprintf(triesUnarchivedPathTemplate, try.ID)
	os.MkdirAll(path, 0666)

	err = ioutil.WriteFile(archivedPath, content, 0666)

	if err != nil {
		return err
	}

	archiver.Zip.Open(archivedPath, unarchivedPath)

	err = t.setCaseStatus(try)

	return err
}

func (t *Tries) setCaseStatus(try *models.Try) error {
	casesRepo := &Cases{}
	sessionsRepo := &Sessions{}

	sessionCase := casesRepo.Find(try.CaseID)
	sessionCase.Status = t.getCaseStatus(try)
	err := casesRepo.Update(sessionCase)

	if err != nil {
		return err
	}

	return sessionsRepo.UpdateSessionStatus(sessionCase.SessionID)
}

func (t *Tries) getCaseStatus(try *models.Try) status.Status {
	var casesTriesCount int

	tryStatus := getTryStatus(try.ExitStatus)

	if tryStatus == status.Success {
		return status.Success
	}

	query := func(db *gorm.DB) {
		db.Where(&models.Try{CaseID: try.CaseID}).Count(&casesTriesCount)
	}

	t.context.Execute(query)

	if casesTriesCount >= config.GetApplicationConfig().MaxRetries {
		return status.Failed
	}
	return status.Pending
}

func getTryStatus(exitCode string) status.Status {
	if exitCode == "0" {
		return status.Success
	}
	return status.Failed
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
