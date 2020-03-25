package repositories

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"malibu-server/data"
	"malibu-server/data/models"
	"malibu-server/data/models/status"
	"malibu-server/utils"

	"github.com/jinzhu/gorm"
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

	utils.Unzip(archivedPath, unarchivedPath)
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

	fmt.Printf("Try status: %s\n", tryStatus.String())

	query := func(db *gorm.DB) {
		db.Model(&models.Try{}).Where(&models.Try{CaseID: try.CaseID}).Count(&casesTriesCount)
	}

	t.context.Execute(query)

	if casesTriesCount >= utils.GetApplicationConfig().MaxRetries {
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

func (t *Tries) FindRawTrySteps(tryID int) []string {
	var result []string

	unarchivedPath := fmt.Sprintf(triesUnarchivedPathTemplate, tryID)
	files, err := ioutil.ReadDir(fmt.Sprintf("./%s/out/", unarchivedPath))

	if err != nil {
		fmt.Println(err.Error())
		return result
	}

	resultMap := make(map[string]bool)

	for _, file := range files {
		filename := path.Base(file.Name())

		if strings.Contains(filename, `SeleniumSessionID`) {
			continue
		}

		extension := filepath.Ext(filename)

		entry := filename[0 : len(filename)-len(extension)]

		if _, value := resultMap[entry]; !value {
			resultMap[entry] = true
			result = append(result, entry)
		}
	}

	return result
}

// Find last successful try by case id
func (t *Tries) FindLastSuccessfulTry(caseID int) *models.Try {
	var try models.Try
	query := func(db *gorm.DB) {
		db.Raw("SELECT * from tries where case_id in ( SELECT id from cases where command_line=(SELECT command_line FROM cases where ID=?)) and exit_status = '0' ORDER BY id desc limit 1", caseID).Scan(&try)
		//		db.First(&try,)
		//		db
	}

	error := t.context.Execute(query)

	if error != nil {
		return nil
	}
	return &try
}
